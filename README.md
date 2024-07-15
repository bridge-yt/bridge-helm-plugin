## Bridge Application: README

### Overview

The Bridge application is designed to manage and register Kubernetes resources created by Helm in a centralized API. It allows for seamless integration and automation by interacting with Helm resources and ensuring they are registered in the Bridge API for centralized configuration management. Additionally, it provides a way to translate placeholders in a `values.yaml` file with actual values from the Bridge API.

### Key Components

1. **HelmResource Struct**:
    - Represents a Kubernetes resource managed by Helm.
    - Fields:
      - `Name`: The name of the resource.
      - `Kind`: The type of the resource (e.g., Pod, Service, Ingress).
      - `Namespace`: The namespace of the resource.
      - `BridgeRegister`: A flag indicating if the resource should be registered in the Bridge API.
      - `Details`: A map containing detailed information about the resource.

2. **Register Function**:
    - Registers Helm resources with the Bridge API.
    - Steps:
      1. Retrieves in-cluster Kubernetes configuration.
      2. Extracts Helm release name and namespace from environment variables.
      3. Retrieves Helm resources using the Helm CLI.
      4. Registers resources in the Bridge API if `BridgeRegister` is true.

3. **Translate Function**:
    - Replaces placeholders in a `values.yaml` file with actual values from the Bridge API.
    - Steps:
      1. Reads the `values.yaml` file.
      2. Extracts placeholders in the format `{{ bridge.<resource>.<field> }}`.
      3. Fetches values from the Bridge API and replaces the placeholders in the `values.yaml` file.

### Usage

#### Environment Setup

1. **Set Environment Variables**:
    - `HELM_RELEASE_NAME`: The name of the Helm release.
    - `HELM_NAMESPACE`: The namespace of the Helm release.

2. **API Endpoint**:
    - Ensure the Bridge API endpoint is specified in the configuration.

#### Example Files

**values.yaml**:
This file contains Helm values and can include placeholders for values that need to be fetched from the Bridge API.

```yaml
serviceAccountName: "{{ bridge.example-resource.arn }}"
```

- `serviceAccountName`: The placeholder `{{ bridge.example-resource.arn }}` will be replaced with the actual ARN of the resource named `example-resource` from the Bridge API.

**templates/pod.yaml**:
This Helm template defines a Pod and includes a label to indicate it should be registered with the Bridge API.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  labels:
    bridgeRegister: "true"
spec:
  containers:
    - name: nginx
      image: nginx
```

- `bridgeRegister: "true"`: This label indicates that the Pod should be registered in the Bridge API.

### Steps to Use

1. **Register Resources**:
    - Ensure your Kubernetes cluster and Helm release are properly configured.
    - Set the environment variables `HELM_RELEASE_NAME` and `HELM_NAMESPACE`.
    - Run the `Register` function to register the resources in the Bridge API.

2. **Translate Placeholders**:
    - Ensure the `values.yaml` file contains placeholders in the format `{{ bridge.<resource>.<field> }}`.
    - Run the `Translate` function to replace the placeholders with actual values from the Bridge API.

### Example Workflow

1. **Register Resources**:
    - Execute the `Register` function to identify and register resources created by a Helm release in the Bridge API.

2. **Translate `values.yaml`**:
    - Run the `Translate` function to replace placeholders in the `values.yaml` file with the actual values from the Bridge API, ensuring the configuration is up-to-date and accurate.

By following this workflow, you can ensure that your Helm-managed Kubernetes resources are centrally registered and their configurations are seamlessly managed and updated across your infrastructure.
