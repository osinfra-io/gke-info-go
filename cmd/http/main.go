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
//         mux.HandleFunc("/gke-info", metadata.GetMetadata)

//         port := "8080"
//         if envPort := os.Getenv("PORT"); envPort != "" {
//                 port = envPort
//         }

//         if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
//                 log.Fatalf("Failed to start server: %v", err)
//         }
// }

func main() {
        mux := http.NewServeMux() // Use standard http.ServeMux
        mux.HandleFunc("/gke-info", metadata.GetMetadata)

        port := "8080"
        if envPort := os.Getenv("PORT"); envPort != "" {
                port = envPort
        }

        if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
                log.Fatalf("Failed to start server: %v", err)
        }
}
