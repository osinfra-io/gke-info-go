package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gke-info/internal/metadata"
)

func TestMain(t *testing.T) {
	// Set up a test server
	mux := http.NewServeMux()
	mux.HandleFunc("/gke-info-go/metadata/", metadata.MetadataHandler)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Set the PORT environment variable to the test server's port
	os.Setenv("PORT", ts.Listener.Addr().String())
	defer os.Unsetenv("PORT")

	// Test cases
	tests := []struct {
		url          string
		expectedCode int
		expectedBody string
		isJSON       bool
	}{
		{ts.URL + "/gke-info-go/metadata/cluster-name", http.StatusOK, `{"cluster-name":"test-cluster-name"}`, true},
		{ts.URL + "/gke-info-go/metadata/cluster-location", http.StatusOK, `{"cluster-location":"test-cluster-location"}`, true},
		{ts.URL + "/gke-info-go/metadata/instance-zone", http.StatusOK, `{"instance-zone":"us-central1-a"}`, true},
		{ts.URL + "/gke-info-go/metadata/unknown", http.StatusBadRequest, "Unknown metadata type\n", false},
	}

	// Mock FetchMetadata for the tests
	originalFetchMetadata := metadata.FetchMetadata
	defer func() { metadata.FetchMetadata = originalFetchMetadata }()
	metadata.FetchMetadata = func(url string) (string, error) {
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

	// Run the test cases
	for _, test := range tests {
		resp, err := http.Get(test.url)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedCode, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if test.isJSON {
			assert.JSONEq(t, test.expectedBody, string(body))
		} else {
			assert.Equal(t, test.expectedBody, string(body))
		}
	}
}
