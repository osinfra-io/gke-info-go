package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gke-info/internal/metadata"
	// "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	// "gopkg.in/DataDog/dd-trace-go.v1/profiler"
	// httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

// func main() {
//         tracer.Start()
//         defer tracer.Stop()

//         err := profiler.Start(
//                 profiler.WithProfileTypes(
//                         profiler.CPUProfile,
//                         profiler.HeapProfile,
//                 ),
//         )
//         if err != nil {
//                 log.Fatal(err)
//         }
//         defer profiler.Stop()

//         mux := httptrace.NewServeMux()
//         mux.HandleFunc("/metadata", metadata.MetadataHandler)

//         port := "8080"
//         if envPort := os.Getenv("PORT"); envPort != "" {
//                 port = envPort
//         }

//         if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
//                 log.Fatalf("Failed to start server: %v", err)
//         }
// }

func main() {
	http.HandleFunc("/metadata", metadata.MetadataHandler)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
        // If successful log the message
        log.Printf("Server started successfully on port %s...\n", port)
}
