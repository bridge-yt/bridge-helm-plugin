package resources

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

type HelmResource struct {
    Name           string            `json:"name"`
    Kind           string            `json:"kind"`
    Namespace      string            `json:"namespace"`
    BridgeRegister bool              `json:"bridgeRegister"`
    Details        map[string]string `json:"details"`
}

func Register(apiURL string) {
    if apiURL == "" {
        log.Fatal("API URL is required")
    }

    config, err := rest.InClusterConfig()
    if err != nil {
        log.Fatalf("Failed to create in-cluster config: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    releaseName := getReleaseName()
    namespace := getNamespace()
    resources, err := getHelmResources(releaseName, namespace)
    if err != nil {
        log.Fatalf("Failed to get Helm resources: %v", err)
    }

    for _, resource := range resources {
        if resource.BridgeRegister {
            err := registerResource(apiURL, resource, clientset)
            if err != nil {
                log.Printf("Failed to register resource %s: %v", resource.Name, err)
            }
        }
    }
}

func getReleaseName() string {
    releaseName := os.Getenv("HELM_RELEASE_NAME")
    if releaseName == "" {
        log.Fatal("HELM_RELEASE_NAME environment variable is required")
    }
    return releaseName
}

func getNamespace() string {
    namespace := os.Getenv("HELM_NAMESPACE")
    if namespace == "" {
        log.Fatal("HELM_NAMESPACE environment variable is required")
    }
    return namespace
}

func getHelmResources(releaseName, namespace string) ([]HelmResource, error) {
    cmd := exec.Command("helm", "get", "manifest", releaseName, "-n", namespace)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("failed to get Helm manifest: %v", err)
    }

    var resources []HelmResource
    var currentResource *HelmResource
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        line = strings.TrimSpace(line)

        if strings.HasPrefix(line, "kind:") {
            if currentResource != nil {
                resources = append(resources, *currentResource)
            }
            currentResource = &HelmResource{
                Kind: strings.TrimPrefix(line, "kind: "),
            }
        } else if strings.HasPrefix(line, "---") {
            currentResource = nil
        } else if currentResource != nil {
            if strings.HasPrefix(line, "name:") {
                currentResource.Name = strings.TrimPrefix(line, "name: ")
            } else if strings.HasPrefix(line, "namespace:") {
                currentResource.Namespace = strings.TrimPrefix(line, "namespace: ")
            } else if strings.Contains(line, "bridgeRegister: true") {
                currentResource.BridgeRegister = true
            }
        }
    }

    if currentResource != nil {
        resources = append(resources, *currentResource)
    }
    return resources, nil
}

func registerResource(apiURL string, resource HelmResource, clientset *kubernetes.Clientset) error {
    var details map[string]string
    switch resource.Kind {
    case "Pod":
        details = getPodDetails(clientset, resource)
    case "Service":
        details = getServiceDetails(clientset, resource)
    case "Ingress":
        details = getIngressDetails(clientset, resource)
    case "Deployment":
        details = getDeploymentDetails(clientset, resource)
    case "StatefulSet":
        details = getStatefulSetDetails(clientset, resource)
    case "ConfigMap":
        details = getConfigMapDetails(clientset, resource)
    case "Secret":
        details = getSecretDetails(clientset, resource)
    default:
        return fmt.Errorf("Unsupported resource type: %s", resource.Kind)
    }

    resource.Details = details

    url := fmt.Sprintf("%s/resource/%s/%s", apiURL, resource.Namespace, resource.Name)
    data := map[string]interface{}{
        "name":          resource.Name,
        "resource_type": resource.Kind,
        "details":       resource.Details,
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return fmt.Errorf("Failed to marshal JSON: %v", err)
    }

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("Failed to create HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("Failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("Failed to register resource: %s", resp.Status)
    }
    return nil
}
