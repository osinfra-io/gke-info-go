package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gke-info/internal/observability"
)

// Metadata server URLs
const (
	ClusterNameURL     = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"
	ClusterLocationURL = "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-location"
	InstanceZoneURL    = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
)

type MetadataFetcher interface {
    FetchMetadata(ctx context.Context, url string) (string, error)
}

func FetchMetadata(ctx context.Context, url string) (string, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        observability.ErrorWithContext(ctx, fmt.Sprintf("Error creating request: %v", err))
        return "", err
    }

    req.Header.Add("Metadata-Flavor", "Google")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        observability.ErrorWithContext(ctx, fmt.Sprintf("Error executing request: %v", err))
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to get metadata from %s, status code: %d, response: %s", url, resp.StatusCode, resp.Status)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        observability.ErrorWithContext(ctx, fmt.Sprintf("Error reading response body: %v", err))
        return "", err
    }

    return string(body), nil
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    if _, err := w.Write([]byte("OK")); err != nil {
        observability.ErrorWithContext(r.Context(), fmt.Sprintf("Error writing response: %v", err))
        http.Error(w, "Failed to write response", http.StatusInternalServerError)
    }
}

func MetadataHandler(fetchMetadataFunc func(ctx context.Context, url string) (string, error)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        observability.InfoWithContext(r.Context(), fmt.Sprintf("Received request for %s", r.URL.Path))

        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            observability.ErrorWithContext(r.Context(), fmt.Sprintf("Invalid request: %s", r.URL.Path))
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
            observability.ErrorWithContext(r.Context(), fmt.Sprintf("Unknown metadata type: %s", metadataType))
            http.Error(w, "Unknown metadata type", http.StatusBadRequest)
            return
        }

        metadata, err := fetchMetadataFunc(r.Context(), url)
        if err != nil {
            observability.ErrorWithContext(r.Context(), fmt.Sprintf("Failed to fetch metadata: %v", err))
            http.Error(w, fmt.Sprintf("Failed to fetch metadata: %v", err), http.StatusInternalServerError)
            return
        }

        if metadataType == "instance-zone" {
            instanceZoneParts := strings.Split(metadata, "/")
            if len(instanceZoneParts) > 0 {
                metadata = instanceZoneParts[len(instanceZoneParts)-1]
            } else {
                observability.ErrorWithContext(r.Context(), fmt.Sprintf("Unexpected format for instance-zone metadata: %s", metadata))
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
}
