package resources

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

func Translate(apiURL string) {
    if apiURL == "" {
        log.Fatal("API URL is required")
    }

    data, err := ioutil.ReadFile("values.yaml")
    if err != nil {
        log.Fatalf("Failed to read values.yaml: %v", err)
    }

    content := string(data)
    placeholders := extractPlaceholders(content)
    for _, placeholder := range placeholders {
        parts := strings.Split(placeholder, ".")
        if len(parts) != 3 {
            continue
        }
        resource := parts[1]
        field := parts[2]
        value, err := fetchValueFromBridge(apiURL, resource, field)
        if err != nil {
            log.Fatalf("Failed to fetch value for %s: %v", placeholder, err)
        }
        content = strings.ReplaceAll(content, fmt.Sprintf("{{ %s }}", placeholder), value)
    }

    err = ioutil.WriteFile("values.yaml", []byte(content), 0644)
    if err != nil {
        log.Fatalf("Failed to write values.yaml: %v", err)
    }
}

func extractPlaceholders(content string) []string {
    var placeholders []string
    for _, line := range strings.Split(content, "\n") {
        if strings.Contains(line, "{{ bridge.") {
            start := strings.Index(line, "{{ bridge.") + 3
            end := strings.Index(line[start:], " }}")
            if end != -1 {
                placeholders = append(placeholders, line[start:start+end])
            }
        }
    }
    return placeholders
}

func fetchValueFromBridge(apiURL, resource, field string) (string, error) {
    url := fmt.Sprintf("%s/resource/%s", apiURL, resource)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Error fetching resource: %s", resp.Status)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    if value, ok := result[field]; ok {
        return fmt.Sprintf("%v", value), nil
    }
    return "", fmt.Errorf("Field %s not found in resource %s", field, resource)
}
