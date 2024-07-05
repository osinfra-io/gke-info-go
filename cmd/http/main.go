package main

import (
	"log"
	"net/http"
	"os"

	"gke-info/internal/metadata"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metadata", metadata.MetadataHandler)
	mux.HandleFunc("/health", metadata.HealthCheckHandler)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
