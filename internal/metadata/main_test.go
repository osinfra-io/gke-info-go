package metadata_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gke-info/internal/metadata"
)

// MockFetchMetadata is a mock implementation of the FetchMetadata function
func MockFetchMetadata(url string) (string, error) {
	switch url {
	case metadata.ClusterNameURL:
		return "test-cluster-name", nil
	case metadata.ClusterLocationURL:
		return "test-cluster-location", nil
	case metadata.InstanceZoneURL:
		return "projects/1234567890/zones/us-central1-a", nil
	default:
		return "", fmt.Errorf("unknown URL: %s", url)
	}
}

// TestFetchMetadata tests the FetchMetadata function
func TestFetchMetadata(t *testing.T) {
	originalFetchMetadata := metadata.FetchMetadata
	defer func() { metadata.FetchMetadata = originalFetchMetadata }()
	metadata.FetchMetadata = MockFetchMetadata

	tests := []struct {
		url      string
		expected string
	}{
		{metadata.ClusterNameURL, "test-cluster-name"},
		{metadata.ClusterLocationURL, "test-cluster-location"},
		{metadata.InstanceZoneURL, "projects/1234567890/zones/us-central1-a"},
	}

	for _, test := range tests {
		result, err := metadata.FetchMetadata(test.url)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, result)
	}
}

// TestMetadataHandler tests the MetadataHandler function
func TestMetadataHandler(t *testing.T) {
	originalFetchMetadata := metadata.FetchMetadata
	defer func() { metadata.FetchMetadata = originalFetchMetadata }()
	metadata.FetchMetadata = MockFetchMetadata

	tests := []struct {
		url          string
		expectedCode int
		expectedBody string
		isJSON       bool
	}{
		{"/metadata/cluster-name", http.StatusOK, `{"cluster-name":"test-cluster-name"}`, true},
		{"/metadata/cluster-location", http.StatusOK, `{"cluster-location":"test-cluster-location"}`, true},
		{"/metadata/instance-zone", http.StatusOK, `{"instance-zone":"us-central1-a"}`, true},
		{"/metadata/unknown", http.StatusBadRequest, "Unknown metadata type", false},
	}

	for _, test := range tests {
		req, err := http.NewRequest("GET", test.url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(metadata.MetadataHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, test.expectedCode, rr.Code)

		body, err := io.ReadAll(rr.Body)
		assert.NoError(t, err)
		if test.isJSON {
			assert.JSONEq(t, test.expectedBody, strings.TrimSpace(string(body)))
		} else {
			assert.Equal(t, test.expectedBody, strings.TrimSpace(string(body)))
		}
	}
}
