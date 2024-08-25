package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gke-info/internal/metadata"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var log = logrus.New()

// main initializes the HTTP server and sets up the routes.
func main() {
	// Set log output to stdout and use JSON formatter
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})

	tracer.Start()
	defer tracer.Stop()

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.WithField("error", err).Warn("Failed to start profiler")
	}
	defer profiler.Stop()

	mux := httptrace.NewServeMux()
	mux.HandleFunc("/gke-info-go/metadata/", metadata.MetadataHandler)
	mux.HandleFunc("/gke-info-go/health", metadata.HealthCheckHandler)

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.WithField("port", port).Info("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithField("error", err).Fatal("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.WithField("error", err).Fatal("Server forced to shutdown")
	}

	log.Info("Server exiting")
}
