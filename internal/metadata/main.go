package metadata

import (
        "fmt"
        "io"
        "net/http"
)

func GetMetadata(w http.ResponseWriter, r *http.Request) {
    // Extract the path parameter from the URL
    metadataType := r.URL.Path[len("/gke-info/"):]

    var metadataURL string

    // Only allow fetching the cluster name for security reasons
    if metadataType == "cluster-name" {
    } else {
        http.Error(w, "Request not allowed", http.StatusForbidden)
        return
    }

    // Construct the metadata URL dynamically
    metadataURL = fmt.Sprintf("http://metadata.google.internal/computeMetadata/v1/instance/attributes/%s", metadataType)

    // Set the metadata request headers
    req, err := http.NewRequest("GET", metadataURL, nil)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
        return
    }
    req.Header.Set("Metadata-Flavor", "Google")

    // Send the request to the metadata server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error retrieving metadata: %v", err), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        http.Error(w, fmt.Sprintf("Error retrieving metadata, status code: %d", resp.StatusCode), http.StatusInternalServerError)
        return
    }

    // Set the Content-Type header to JSON
    w.Header().Set("Content-Type", "application/json")

    // Copy the response from the metadata server to the output
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
        return
    }
}
