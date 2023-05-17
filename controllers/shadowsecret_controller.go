/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"io"
	"strings"

	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/integration/gcpkms"
	"github.com/google/tink/go/keyset"
	corev1alpha1 "github.com/phoban01/rune/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ShadowSecretReconciler reconciles a ShadowSecret object
type ShadowSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.rune.io,resources=shadowsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.rune.io,resources=shadowsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.rune.io,resources=shadowsecrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop
func (r *ShadowSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// get the shadow secret object
	obj := &corev1alpha1.ShadowSecret{}
	if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get object: %w", err)
	}
	logger.Info("got shadow secret")

	// get the secret store object
	secretStore := &corev1alpha1.SecretStore{}
	if err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: obj.GetNamespace(),
		Name:      obj.Spec.SecretStoreRef.Name,
	}, secretStore); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get object: %w", err)
	}
	logger.Info("got secret store")

	// TODO: validate that the service account can get the secret at path

	// get the registry secret
	regCred := &corev1.Secret{}
	if err := r.Client.Get(ctx,
		types.NamespacedName{
			Namespace: secretStore.GetNamespace(),
			Name:      secretStore.Spec.Registry.SecretRef.Name}, regCred); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get object: %w", err)
	}

	// get the kms secret
	kmsCred := &corev1.Secret{}
	if err := r.Client.Get(ctx,
		types.NamespacedName{
			Namespace: secretStore.GetNamespace(),
			Name:      secretStore.Spec.KMS.SecretRef.Name}, kmsCred); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get object: %w", err)
	}

	// auth to kms
	gcpClient, err := gcpkms.NewClientWithOptions(ctx, secretStore.Spec.KMS.Value, option.WithCredentialsJSON(kmsCred.Data["credentials.json"]))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create gcp client: %w", err)
	}
	registry.RegisterKMSClient(gcpClient)

	dek := aead.AES128CTRHMACSHA256KeyTemplate()
	kh, err := keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(secretStore.Spec.KMS.Value, dek))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create key template: %w", err)
	}

	a, err := aead.New(kh)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get object: %w", err)
	}

	// auth to registry
	authConfig := authn.AuthConfig{Auth: string(regCred.Data[".dockerconfigjson"])}

	// fetch the cipher text from registry
	parseRef, err := name.ParseReference(fmt.Sprintf("%s/rune/%s:%s", strings.TrimPrefix(secretStore.Spec.Registry.URL, "oci://"), obj.Spec.Path, obj.Spec.Version))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	descriptor, err := remote.Get(parseRef, remote.WithAuth(authn.FromConfig(authConfig)))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	// decrypt using tink
	image, err := descriptor.Image()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	layers, err := image.Layers()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	aad := []byte("this data needs to be authenticated, but not encrypted")
	data, err := layers[0].Uncompressed()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	ct, err := io.ReadAll(data)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	pt, err := a.Decrypt(ct, aad)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get auth: %w", err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: obj.GetNamespace(),
			Name:      obj.GetName(),
		},
	}
	// generate a secret
	controllerutil.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if secret.ObjectMeta.CreationTimestamp.IsZero() {
			if err := controllerutil.SetOwnerReference(obj, secret, r.Scheme); err != nil {
				return fmt.Errorf("failed to set owner reference: %w", err)
			}
		}
		secret.StringData = map[string]string{
			"test": string(pt),
		}
		return nil
	})

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShadowSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.ShadowSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
