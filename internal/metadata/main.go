package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Metadata server URLs
const (
	clusterNameURL     = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"
	clusterLocationURL = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-location"
	instanceZoneURL    = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
)

// fetchMetadata fetches metadata from the provided URL
func fetchMetadata(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Adding the metadata flavor header as required by GCP metadata server
	req.Header.Add("Metadata-Flavor", "Google")

	client := &http.Client{}
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

// metadataHandler handles the /metadata/* endpoint
func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	metadataType := pathParts[2]
	var url string

	switch metadataType {
	case "cluster-name":
		url = clusterNameURL
	case "cluster-location":
		url = clusterLocationURL
	case "instance-zone":
		url = instanceZoneURL
	default:
		http.Error(w, "Unknown metadata type", http.StatusBadRequest)
		return
	}

	metadata, err := fetchMetadata(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If instance zone, extract just the zone from the full zone URL
	if metadataType == "instance-zone" {
		instanceZoneParts := strings.Split(metadata, "/")
		if len(instanceZoneParts) > 0 {
			metadata = instanceZoneParts[len(instanceZoneParts)-1]
		}
	}

	response := map[string]string{metadataType: metadata}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
