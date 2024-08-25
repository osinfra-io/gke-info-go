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

	dd_logrus "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var log = logrus.New()

// main initializes the HTTP server and sets up the routes.
func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	// Add Datadog context log hook
	logrus.AddHook(&dd_logrus.DDContextLogHook{})

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

	// Create a context
	ctx := context.Background()

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
		logrus.WithContext(ctx).Info("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithContext(ctx).Fatal("Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.WithContext(ctx).Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logrus.WithContext(ctx).Fatal("Server forced to shutdown")
	}

	logrus.WithContext(ctx).Info("Server exiting")
}
