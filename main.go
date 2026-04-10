package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

// Service metadata
var (
	serviceName = "go-restapi-server"
	version     = "dev" // Default fallback version
)

var people []Person

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Person not found"}`))
}
func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	var person Person
	err := json.NewDecoder(req.Body).Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Invalid JSON"}`))
		return
	}
	
	// Validate required fields
	if person.Firstname == "" || person.Lastname == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Firstname and Lastname are required"}`))
		return
	}
	
	// Generate ID for new person
	person.ID = generateID()
	people = append(people, person)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// Helper function to generate a new ID
func generateID() string {
	maxID := 0
	for _, p := range people {
		// Parse existing ID as integer to find max
		var id int
		_, err := fmt.Sscanf(p.ID, "%d", &id)
		if err == nil && id > maxID {
			maxID = id
		}
	}
	return fmt.Sprintf("%d", maxID+1)
}
func UpdatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			// Decode the new person data
			var person Person
			err := json.NewDecoder(req.Body).Decode(&person)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Invalid JSON"}`))
				return
			}

			// Validate required fields
			if person.Firstname == "" || person.Lastname == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"Firstname and Lastname are required"}`))
				return
			}

			// Set the ID to the original ID
			person.ID = params["id"]

			// Update the person in place
			people[index] = person

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(person)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Person not found"}`))
}
func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			// Remove the person from the slice
			people = append(people[:index], people[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Person not found"}`))
}

func HealthEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// VersionEndpoint returns service metadata
func VersionEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"service": serviceName,
		"version": version,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "Ernest", Lastname: "Hemingway", Address: &Address{City: "Dublin", State: "CA"}})
	people = append(people, Person{ID: "2", Firstname: "George", Lastname: "Orwell"})
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", UpdatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	router.HandleFunc("/health", HealthEndpoint).Methods("GET")
	router.HandleFunc("/version", VersionEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))
}
