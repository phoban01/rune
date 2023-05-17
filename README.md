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
    url: oci://my-oci-registry.example.com
    secretRef:
      name: my-registry-credentials
  kms:
    provider: google
    value:
    secretRef:
      name: my-kms-credentials
```

### **OCI Registry as Secret Store**

OCI registries offer a range of advantages that make them an excellent choice for storing encrypted data. These advantages include secure signing, efficient mirroring, and broad accessibility, positioning OCI registries as a robust and versatile option for encrypted data storage.

Secure signing is a crucial factor when it comes to encrypted data storage, and OCI registries excel in this aspect. OCI registries support cryptographic signing of container images and artifacts using digital signatures. This signing mechanism ensures the integrity and authenticity of the encrypted data stored in the registry. By verifying the digital signatures, organizations can trust the source and integrity of the encrypted data, providing an additional layer of security and mitigating the risks of tampering or unauthorized modifications.

Efficient mirroring is another notable advantage of OCI registries for encrypted data storage. OCI registries offer the capability to create mirrors or replicas of the registry across multiple instances or locations. This mirroring mechanism enables organizations to achieve data redundancy and high availability, ensuring uninterrupted access to the encrypted data. Mirroring also improves the overall performance by enabling localized access to the encrypted data, reducing latency and network congestion. With efficient mirroring capabilities, OCI registries provide organizations with a reliable and scalable solution for storing and accessing encrypted data.

Moreover, OCI registries are widely used and accessible in various environments. They have gained significant adoption across the container ecosystem, becoming a standard for storing all sorts of data, including encrypted data. OCI registries are compatible with different container runtimes, orchestrators, and deployment platforms, making them accessible in diverse computing environments. This wide acceptance and accessibility ensure that encrypted data stored in OCI registries can be seamlessly utilized and integrated into different applications, services, and systems across different infrastructures.

###**Getting Started**

To get started with Rune, refer to the documentation for installation instructions, usage guides, and configuration details. The documentation provides step-by-step instructions, sample YAML definitions, and CLI usage examples to assist users in effectively utilizing the Rune service for their secret management needs.


**Community and Contributions**

Rune is an open-source project, and we welcome contributions from the community. If you encounter issues, have suggestions, or would like to contribute to the project, please visit the Rune GitHub repository. We value the input and participation of the community in improving the features, functionality, and security of Rune.

**License**

Rune is released under the Apache License 2.0. Please refer to the LICENSE file for more details.

**Acknowledgments**

Rune is built upon the efforts of various open-source projects and libraries. We extend our gratitude to the contributors and maintainers of these projects for their valuable work and contributions.
