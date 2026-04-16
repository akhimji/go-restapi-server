package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go-restapi-server/internal/handlers"
	"go-restapi-server/internal/observability"
	"go-restapi-server/internal/store"
)

// Service metadata
var (
	serviceName = "go-restapi-server"
	version     = "dev" // Default fallback version
)

// VersionEndpoint returns the version of the service
func VersionEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"service": serviceName, "version": version})
}

func main() {
	// Initialize slog with JSON handler for structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	router := mux.NewRouter()

	// Create store
	personStore := store.NewInMemoryPersonStore()

	// Create handlers
	peopleHandler := handlers.NewPeopleHandler(personStore)

	// Register routes
	router.HandleFunc("/people", peopleHandler.GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people", peopleHandler.CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", peopleHandler.GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", peopleHandler.UpdatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/people/{id}", peopleHandler.PatchPersonEndpoint).Methods("PATCH")
	router.HandleFunc("/people/{id}", peopleHandler.DeletePersonEndpoint).Methods("DELETE")
	router.HandleFunc("/health", handlers.HealthEndpoint).Methods("GET")
	router.HandleFunc("/version", VersionEndpoint).Methods("GET")

	// Wrap router with request ID middleware first, then logging middleware
	router.Use(observability.RequestIDMiddleware)
	router.Use(observability.LoggingMiddleware)

	// Create server with timeout
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server exited gracefully")
}