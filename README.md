## **Rune - Cloud Native Secrets Management**

Rune is an open-source solution designed to address the challenges of secret management in Kubernetes environments. It provides a secure and seamless approach to storing, retrieving, and utilizing secrets while ensuring strong access control, authentication, and encryption mechanisms.

Secrets are stored encrypted in an OCI (Open Container Initiative) registry, ensuring their confidentiality and integrity. Rune leverages Google Tink, an open-source cryptographic library, to ensure robust encryption and decryption of secrets.  The `rune` CLI is used to read, write, encrypt, and decrypt secrets, while the `rune-controller` makes secrets accessible within the Kubernetes cluster by decrypting them using any method supported by Tink. Secret integrity and provenance is verified `cosign`.

To enable a registry, a `RuneStore` CRD is created, providing the access credentials. Rune supports path-based RBAC (Role-Based Access Control), ensuring that only authorized entities can access specific secrets based on a policy which can be stored in the registry itself.

Kubernetes workloads can retrieve secrets in two ways:

- **Rune Secret CRD**: The `Rune Secret` CRD specifies the interval, `RuneStore` reference, service account name, and path. The `rune-controller` uses the service account as the principal when validating RBAC.

- **Rune API**: Workloads can request secrets from the `rune-controller` server at runtime using the HTTP/gRPC API. The API requires a service account JWT for RBAC validation, and a Go client SDK is provided for easy integration.

By combining cloud-native principles with robust encryption and access control mechanisms, Rune simplifies secret management in Kubernetes environments while ensuring the highest level of security and confidentiality for sensitive information.

### **OCI Registry as Secret Store**
OCI registries offer a range of advantages that make them an excellent choice for storing encrypted data. These advantages include secure signing, efficient mirroring, and broad accessibility, positioning OCI registries as a robust and versatile option for encrypted data storage.

Secure signing is a crucial factor when it comes to encrypted data storage, and OCI registries excel in this aspect. OCI registries support cryptographic signing of container images and artifacts using digital signatures. This signing mechanism ensures the integrity and authenticity of the encrypted data stored in the registry. By verifying the digital signatures, organizations can trust the source and integrity of the encrypted data, providing an additional layer of security and mitigating the risks of tampering or unauthorized modifications.

Efficient mirroring is another notable advantage of OCI registries for encrypted data storage. OCI registries offer the capability to create mirrors or replicas of the registry across multiple instances or locations. This mirroring mechanism enables organizations to achieve data redundancy and high availability, ensuring uninterrupted access to the encrypted data. Mirroring also improves the overall performance by enabling localized access to the encrypted data, reducing latency
and network congestion. With efficient mirroring capabilities, OCI registries provide organizations with a reliable and scalable solution for storing and accessing encrypted data.

Moreover, OCI registries are widely used and accessible in various environments. They have gained significant adoption across the container ecosystem, becoming a standard for storing all sorts of data, including encrypted data. OCI registries are compatible with different container runtimes, orchestrators, and deployment platforms, making them accessible in diverse computing environments. This wide acceptance and accessibility ensure that encrypted data stored in OCI registries can be seamlessly utilized and integrated into different applications, services, and systems across different infrastructures.

### **Google Tink for Encryption**
To ensure robust encryption and decryption of secrets, Rune leverages Google Tink, an open-source cryptographic library developed by Google. Tink provides a comprehensive set of cryptographic primitives and high-level APIs, making it easy to implement secure encryption and decryption mechanisms.

When storing secrets in the OCI registry, Rune utilizes Google Tink to encrypt the sensitive information before it is stored. This ensures that the secrets remain confidential and protected, even if unauthorized access to the registry occurs. The encryption process utilizes strong encryption algorithms and best practices to safeguard the secrets.

During secret retrieval, the `rune-controller` uses Google Tink to decrypt the encrypted secrets. The necessary decryption keys and algorithms are securely stored within the RuneStore CRD. This approach ensures that only authorized entities with the appropriate credentials can access and decrypt the secrets, maintaining the confidentiality and integrity of the sensitive information.

By utilizing Google Tink for encryption and decryption, Rune benefits from a battle-tested and highly secure cryptographic library. Google Tink follows rigorous security practices, undergoes regular security audits, and incorporates the latest advancements in cryptography. This helps ensure that the secrets stored and retrieved by Rune remain well-protected and resistant to various cryptographic attacks.

Additionally, Google Tink provides support for a wide range of encryption algorithms, key management systems, and cryptographic operations, allowing Rune to adapt and evolve its encryption capabilities based on the specific requirements and preferences of users and organizations. The flexibility and extensibility of Google Tink make it a reliable choice for handling encryption within the Rune secret management solution.

Integrating Google Tink into Rune adds an extra layer of security and trust to the secret management process, enhancing the overall confidentiality and protection of sensitive information within Kubernetes environments.

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
