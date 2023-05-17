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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretStoreSpec defines the desired state of SecretStore
type SecretStoreSpec struct {
	// +required
	Registry RegistrySpec `json:"registry"`

	// +required
	KMS KMSSpec `json:"kms"`
}

// RegistrySpec defines the secret store registry
type RegistrySpec struct {
	// +required
	URL string `json:"url"`

	// +required
	SecretRef corev1.LocalObjectReference `json:"secretRef"`
}

// KMSSpec defines the kms to use
type KMSSpec struct {
	// +required
	Provider string `json:"provider"`

	// +required
	Value string `json:"value"`

	// +required
	SecretRef corev1.LocalObjectReference `json:"secretRef"`
}

// SecretStoreStatus defines the observed state of SecretStore
type SecretStoreStatus struct {
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SecretStore is the Schema for the secretstores API
type SecretStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretStoreSpec   `json:"spec,omitempty"`
	Status SecretStoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SecretStoreList contains a list of SecretStore
type SecretStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretStore{}, &SecretStoreList{})
}
