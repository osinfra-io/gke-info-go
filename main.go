package main

import (
	"fmt"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
	"io"
	"log"
	"net/http"
	"os"
)

func getClusterInfo(w http.ResponseWriter, _ *http.Request, client *http.Client) {
	// GCP metadata URL to retrieve the GKE cluster name
	metadataURL := "http://metadata.google.internal/computeMetadata/v1/instance/attributes/cluster-name"

	// Set the metadata request headers
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Metadata-Flavor", "Google")

	// Send the request to the metadata server
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

func main() {
	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,

			// The profiles below are disabled by
			// default to keep overhead low, but
			// can be enabled as needed.
			// profiler.BlockProfile,
			// profiler.MutexProfile,
			// profiler.GoroutineProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	http.HandleFunc("/gke-info", func(w http.ResponseWriter, r *http.Request) {
		getClusterInfo(w, r, nil)
	})

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
