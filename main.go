package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go-restapi-server/internal/handlers"
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
	router.HandleFunc("/people/{id}", peopleHandler.DeletePersonEndpoint).Methods("DELETE")
	router.HandleFunc("/health", handlers.HealthEndpoint).Methods("GET")
	router.HandleFunc("/version", VersionEndpoint).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}