package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Metadata server URLs
const (
	ClusterNameURL     = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"
	ClusterLocationURL = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-location"
	InstanceZoneURL    = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
)

// fetchMetadata fetches metadata from the provided URL
var FetchMetadata = func(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Adding the metadata flavor header as required by GCP metadata server
	req.Header.Add("Metadata-Flavor", "Google")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get metadata from %s, status code: %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    if _, err := w.Write([]byte("OK")); err != nil {
        http.Error(w, "Failed to write response", http.StatusInternalServerError)
    }
}

// metadataHandler handles the /metadata/* endpoint
func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request: expected /metadata/{type}", http.StatusBadRequest)
		return
	}

	metadataType := pathParts[2]
	var url string

	switch metadataType {
	case "cluster-name":
		url = ClusterNameURL
	case "cluster-location":
		url = ClusterLocationURL
	case "instance-zone":
		url = InstanceZoneURL
	default:
		http.Error(w, "Unknown metadata type", http.StatusBadRequest)
		return
	}

	metadata, err := FetchMetadata(url)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	if metadataType == "instance-zone" {
		instanceZoneParts := strings.Split(metadata, "/")
		if len(instanceZoneParts) > 0 {
			metadata = instanceZoneParts[len(instanceZoneParts)-1]
		}
	}

	response := map[string]string{metadataType: metadata}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
