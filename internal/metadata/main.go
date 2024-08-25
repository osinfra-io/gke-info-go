package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Metadata server URLs
const (
	ClusterNameURL     = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"
	ClusterLocationURL = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-location"
	InstanceZoneURL    = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
)

var log = logrus.New()

func init() {
    // Set log output to stdout and use JSON formatter
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// FetchMetadata fetches metadata from the provided URL and returns it as a string
var FetchMetadata = func(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
		}).Error("Error creating request")
		return "", err
	}

	// Adding the metadata flavor header as required by GCP metadata server
	req.Header.Add("Metadata-Flavor", "Google")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
		}).Error("Error executing request")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get metadata from %s, status code: %d, response: %s", url, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
			"url":   url,
		}).Error("Error reading response body")
		return "", err
	}

	return string(body), nil
}

// HealthCheckHandler responds with a simple "OK" to indicate the service is healthy.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.WithField("error", err).Error("Error writing response")
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// MetadataHandler handles the /gke-info-go/metadata/* endpoint and fetches the requested metadata.
func MetadataHandler(w http.ResponseWriter, r *http.Request) {
	log.WithField("path", r.URL.Path).Info("Received request")

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.WithField("path", r.URL.Path).Error("Invalid request")
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
		log.WithField("metadataType", metadataType).Error("Unknown metadata type")
		http.Error(w, "Unknown metadata type", http.StatusBadRequest)
		return
	}

	metadata, err := FetchMetadata(url)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error":        err,
			"metadataType": metadataType,
			"url":          url,
		}).Error("Failed to fetch metadata")
		http.Error(w, fmt.Sprintf("Failed to fetch metadata: %v", err), http.StatusInternalServerError)
		return
	}

	// Handle special case for instance-zone metadata and write the response.
	if metadataType == "instance-zone" {
		instanceZoneParts := strings.Split(metadata, "/")
		if len(instanceZoneParts) > 0 {
			metadata = instanceZoneParts[len(instanceZoneParts)-1]
		} else {
			log.WithField("metadata", metadata).Error("Unexpected format for instance-zone metadata")
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
