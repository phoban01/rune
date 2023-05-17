# Rune

*Problem Statement: Secure and Efficient Secret Management in Kubernetes*

Managing secrets in a Kubernetes environment can be a challenging task. Traditional approaches often involve manual handling of sensitive information, leading to potential security vulnerabilities, operational inefficiencies, and difficulties in maintaining secret integrity across various applications and microservices. Furthermore, securely retrieving and decrypting secrets while ensuring proper access control and authentication can add complexity to the overall system.

## **Introduction: Rune - Secure and Seamless Secret Retrieval for Kubernetes**

Rune is an open-source solution designed to address the challenges of secret management in Kubernetes environments. It provides a secure and seamless approach to retrieving, decrypting, and utilizing secrets while ensuring strong access control and authentication mechanisms.

Rune stores secrets encrypted in an OCI (Open Container Initiative) registry, ensuring their confidentiality and integrity. The `rune` CLI is used to read, write, encrypt, and decrypt secrets. The `rune-controller` makes secrets accessible within the Kubernetes cluster by decrypting them using the access credentials stored in the RuneStore CRD.

To enable a registry, a `RuneStore` CRD is created, providing the access credentials. Path-based RBAC is enforced on secrets, with the policy stored in the registry and written in CUE.

Users can retrieve secrets in two ways:

- **Rune Secret CRD**: The `Rune Secret` CRD specifies the interval, `RuneStore` reference, service account name, and path. The `rune-controller` uses the service account as the principal when validating RBAC.

- **Rune API**: Workloads can request secrets from the `rune-controller` server at runtime using the HTTP/gRPC API. The API requires a service account JWT for RBAC validation, and a Go client SDK is provided for easy integration.

## **Examples**

1. Rune Secret CRD:

```yaml
apiVersion: core.rune.io/v1alpha1
kind: ShadowSecret
metadata:
  name: my-secret
  namespace: default
spec:
  interval: 1h
  runeStoreRef: my-rune-store
  serviceAccountName: app-team
  path: production/db/postgres
```

2. RuneStore CRD:

```yaml
apiVersion: core.rune.io/v1alpha1
kind: SecretStore
metadata:
  name: my-rune-store
  namespace: default
spec:
  registry:
    url: https://my-oci-registry.example.com
    credentials:
      secretRef:
        name: my-registry-credentials
```

**Getting Started**

To get started with Rune, refer to the documentation for installation instructions, usage guides, and configuration details. The documentation provides step-by-step instructions, sample YAML definitions, and CLI usage examples to assist users in effectively utilizing the Rune service for their secret management needs.

**Community and Contributions**

Rune is an open-source project, and we welcome contributions from the community. If you encounter issues, have suggestions, or would like to contribute to the project, please visit the Rune GitHub repository. We value the input and participation of the community in improving the features, functionality, and security of Rune.

**License**

Rune is released under the Apache License 2.0. Please refer to the LICENSE file for more details.

**Acknowledgments**

Rune is built upon the efforts of various open-source projects and libraries. We extend our gratitude to the contributors and maintainers of these projects for their valuable work and contributions.
