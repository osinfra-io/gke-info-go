package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"istio-test/internal/metadata"
)

func TestMain(t *testing.T) {
    mux := http.NewServeMux()
    mux.HandleFunc("/istio-testst/metadata/", metadata.MetadataHandler(mockFetchMetadata))

    ts := httptest.NewServer(mux)
    defer ts.Close()

    os.Setenv("PORT", ts.Listener.Addr().String())
    defer os.Unsetenv("PORT")

    tests := []struct {
        url          string
        expectedCode int
        expectedBody string
        isJSON       bool
    }{
        {ts.URL + "/istio-testst/metadata/cluster-name", http.StatusOK, `{"cluster-name":"test-cluster-name"}`, true},
        {ts.URL + "/istio-testst/metadata/cluster-location", http.StatusOK, `{"cluster-location":"test-cluster-location"}`, true},
        {ts.URL + "/istio-testst/metadata/instance-zone", http.StatusOK, `{"instance-zone":"us-central1-a"}`, true},
        {ts.URL + "/istio-testst/metadata/unknown", http.StatusBadRequest, "Unknown metadata type\n", false},
    }

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

func mockFetchMetadata(ctx context.Context, url string) (string, error) {
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
