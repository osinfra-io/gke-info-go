package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"istio-test/internal/metadata"
	"istio-test/internal/observability"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

func main() {
	ctx := context.Background()

	observability.Init()

	observability.InfoWithContext(ctx, "Application is starting")

	tracer.Start(tracer.WithRuntimeMetrics())
	defer tracer.Stop()

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		observability.ErrorWithContext(ctx, fmt.Sprintf("Warning: Failed to start profiler: %v", err))
	}
	defer profiler.Stop()

	mux := httptrace.NewServeMux()
	mux.HandleFunc("/istio-test/metadata/", metadata.MetadataHandler(metadata.FetchMetadata))
	mux.HandleFunc("/istio-test/health", metadata.HealthCheckHandler)
	mux.HandleFunc("/", metadata.NotFoundHandler)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := &http.Server{
		Addr:        ":" + port,
		ReadTimeout: 5 * time.Second,
		Handler:     mux,
	}

	go func() {
		observability.InfoWithContext(ctx, fmt.Sprintf("Starting server on port %s...", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			observability.ErrorWithContext(ctx, fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	observability.InfoWithContext(ctx, "Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		observability.ErrorWithContext(ctx, fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	observability.InfoWithContext(ctx, "Server exiting")
}
