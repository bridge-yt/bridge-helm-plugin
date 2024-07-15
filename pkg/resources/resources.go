package resources

import (
    "context"
    "fmt"
    "log"
    "strings"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    corev1 "k8s.io/api/core/v1"
)

func getPodDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    pod, err := clientset.CoreV1().Pods(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get Pod details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace": resource.Namespace,
        "podIP":     pod.Status.PodIP,
        "nodeName":  pod.Spec.NodeName,
        // Add more Pod-specific details as needed
    }
}

func getServiceDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    service, err := clientset.CoreV1().Services(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get Service details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace": resource.Namespace,
        "clusterIP": service.Spec.ClusterIP,
        "ports":     strings.Join(getServicePorts(service.Spec.Ports), ","),
        // Add more Service-specific details as needed
    }
}

func getServicePorts(ports []corev1.ServicePort) []string {
    var portStrings []string
    for _, port := range ports {
        portStrings = append(portStrings, fmt.Sprintf("%d/%s", port.Port, port.Protocol))
    }
    return portStrings
}

func getIngressDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    ingress, err := clientset.NetworkingV1().Ingresses(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get Ingress details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace": resource.Namespace,
        "host":      ingress.Spec.Rules[0].Host,
        // Add more Ingress-specific details as needed
    }
}

func getDeploymentDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    deployment, err := clientset.AppsV1().Deployments(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get Deployment details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace":          resource.Namespace,
        "replicas":           fmt.Sprintf("%d", *deployment.Spec.Replicas),
        "availableReplicas":  fmt.Sprintf("%d", deployment.Status.AvailableReplicas),
        "updatedReplicas":    fmt.Sprintf("%d", deployment.Status.UpdatedReplicas),
        // Add more Deployment-specific details as needed
    }
}

func getStatefulSetDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    statefulSet, err := clientset.AppsV1().StatefulSets(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get StatefulSet details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace":        resource.Namespace,
        "replicas":         fmt.Sprintf("%d", *statefulSet.Spec.Replicas),
        "readyReplicas":    fmt.Sprintf("%d", statefulSet.Status.ReadyReplicas),
        "currentReplicas":  fmt.Sprintf("%d", statefulSet.Status.CurrentReplicas),
        // Add more StatefulSet-specific details as needed
    }
}

func getConfigMapDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    configMap, err := clientset.CoreV1().ConfigMaps(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get ConfigMap details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace": resource.Namespace,
        "dataKeys":  strings.Join(getKeys(configMap.Data), ","),
        // Add more ConfigMap-specific details as needed
    }
}

func getSecretDetails(clientset *kubernetes.Clientset, resource HelmResource) map[string]string {
    secret, err := clientset.CoreV1().Secrets(resource.Namespace).Get(context.TODO(), resource.Name, metav1.GetOptions{})
    if err != nil {
        log.Printf("Failed to get Secret details for %s/%s: %v", resource.Namespace, resource.Name, err)
        return nil
    }
    return map[string]string{
        "namespace": resource.Namespace,
        "dataKeys":  strings.Join(getKeysByte(secret.Data), ","),
        "type":      string(secret.Type),
        // Add more Secret-specific details as needed
    }
}

// Helper function to get keys from a map
func getKeys(data map[string]string) []string {
    keys := make([]string, 0, len(data))
    for k := range data {
        keys = append(keys, k)
    }
    return keys
}

// Helper function to get keys from a map of []byte
func getKeysByte(data map[string][]byte) []string {
    keys := make([]string, 0, len(data))
    for k := range data {
        keys = append(keys, k)
    }
    return keys
}
