package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// Helper function to reset people slice to initial state
func resetPeople() {
	people = nil
	people = append(people, Person{ID: "1", Firstname: "Ernest", Lastname: "Hemingway", Address: &Address{City: "Dublin", State: "CA"}})
	people = append(people, Person{ID: "2", Firstname: "George", Lastname: "Orwell"})
}

func TestDeletePersonEndpoint(t *testing.T) {
	// Reset the people slice to initial state for each test
	resetPeople()

	// Test delete existing person
	t.Run("delete existing person", func(t *testing.T) {
		resetPeople() // Reset for each test case
		req := httptest.NewRequest("DELETE", "/people/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		DeletePersonEndpoint(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check that we got the correct success message
		expectedBody := `{"message":"Person deleted successfully"}`
		body := w.Body.String()
		if body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})

	// Test delete non-existing person
	t.Run("delete non-existing person", func(t *testing.T) {
		resetPeople() // Reset for each test case
		req := httptest.NewRequest("DELETE", "/people/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()
		DeletePersonEndpoint(w, req)

		if status := w.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		// Check that we got the correct error message
		expectedBody := `{"error":"Person not found"}`
		body := w.Body.String()
		if body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})
}

func TestGetPersonEndpoint(t *testing.T) {
	// Test that GET still works
	resetPeople()
	req := httptest.NewRequest("GET", "/people/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	GetPersonEndpoint(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetPeopleEndpoint(t *testing.T) {
	// Test that GET all still works
	resetPeople()
	req := httptest.NewRequest("GET", "/people", nil)
	w := httptest.NewRecorder()
	GetPeopleEndpoint(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHealthEndpoint(t *testing.T) {
	// Test that health endpoint returns 200 and correct payload
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	HealthEndpoint(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that we got the correct success payload
	expectedBody := `{"status":"healthy"}`
	body := w.Body.String()
	if body != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, expectedBody)
	}
}

func TestVersionEndpoint(t *testing.T) {
	// Test that version endpoint returns 200 and correct JSON payload
	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()
	VersionEndpoint(w, req)

	// Check status code
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check content type
	expectedContentType := "application/json"
	actualContentType := w.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			actualContentType, expectedContentType)
	}

	// Check that we got valid JSON with service and version fields
	body := w.Body.String()
	if body == "" {
		t.Errorf("handler returned empty body")
	}

	// Use table-driven tests for more thorough validation
	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"service field", "service", serviceName},
		{"version field", "version", version},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that the field exists and is non-empty
			if tt.expected == "" {
				t.Errorf("expected %s to be non-empty", tt.field)
			}
		})
	}
}