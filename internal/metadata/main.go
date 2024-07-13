package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// FetchMetadata fetches metadata from the provided URL and returns it as a string
var FetchMetadata = func(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "", err
	}

	// Adding the metadata flavor header as required by GCP metadata server
	req.Header.Add("Metadata-Flavor", "Google")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error executing request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get metadata from %s, status code: %d, response: %s", url, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", err
	}

	return string(body), nil
}

// HealthCheckHandler responds with a simple "OK" to indicate the service is healthy.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// MetadataHandler handles the /metadata/* endpoint and fetches the requested metadata.
func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for %s", r.URL.Path)

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.Printf("Invalid request: %s", r.URL.Path)
		http.Error(w, "Invalid request: expected /gke-info-go/metadata/{type}", http.StatusBadRequest)
		return
	}

	metadataType := pathParts[3]
	var url string

	switch metadataType {
	case "cluster-name":
		url = ClusterNameURL
	case "cluster-location":
		url = ClusterLocationURL
	case "instance-zone":
		url = InstanceZoneURL
	default:
		log.Printf("Unknown metadata type: %s", metadataType)
		http.Error(w, "Unknown metadata type", http.StatusBadRequest)
		return
	}

	metadata, err := FetchMetadata(url)
	if err != nil {
		log.Printf("Failed to fetch metadata: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch metadata: %v", err), http.StatusInternalServerError)
		return
	}

	// Handle special case for instance-zone metadata and write the response.
	if metadataType == "instance-zone" {
		instanceZoneParts := strings.Split(metadata, "/")
		if len(instanceZoneParts) > 0 {
			metadata = instanceZoneParts[len(instanceZoneParts)-1]
		} else {
			log.Printf("Unexpected format for instance-zone metadata: %s", metadata)
			http.Error(w, "Unexpected format for instance-zone metadata", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]string{metadataType: metadata}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
