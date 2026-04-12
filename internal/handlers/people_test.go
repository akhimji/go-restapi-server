package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"go-restapi-server/internal/models"
	"go-restapi-server/internal/store"
)

func TestGetPersonEndpoint(t *testing.T) {
	// Create a test store with some data
	store := store.NewInMemoryPersonStore()

	// Add a test person
	person := &models.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	store.Create(person)

	// Create a people handler
	handler := NewPeopleHandler(store)

	// Test case 1: Get existing person
	t.Run("get existing person", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.GetPersonEndpoint(w, req)

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

		// Check response body
		var responsePerson models.Person
		if err := json.NewDecoder(w.Body).Decode(&responsePerson); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		if responsePerson.ID != person.ID {
			t.Errorf("handler returned wrong person ID: got %v want %v",
				responsePerson.ID, person.ID)
		}
		if responsePerson.FirstName != person.FirstName {
			t.Errorf("handler returned wrong first name: got %v want %v",
				responsePerson.FirstName, person.FirstName)
		}
		if responsePerson.LastName != person.LastName {
			t.Errorf("handler returned wrong last name: got %v want %v",
				responsePerson.LastName, person.LastName)
		}
		if responsePerson.Email != person.Email {
			t.Errorf("handler returned wrong email: got %v want %v",
				responsePerson.Email, person.Email)
		}
		if responsePerson.Age != person.Age {
			t.Errorf("handler returned wrong age: got %v want %v",
				responsePerson.Age, person.Age)
		}
	})

	// Test case 2: Get non-existing person
	t.Run("get non-existing person", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()
		handler.GetPersonEndpoint(w, req)

		if status := w.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}

		// Check content type
		expectedContentType := "application/json"
		actualContentType := w.Header().Get("Content-Type")
		if actualContentType != expectedContentType {
			t.Errorf("handler returned wrong content type: got %v want %v",
				actualContentType, expectedContentType)
		}

		// Check error message in response
		expectedBody := `{"error":"Person not found"}`
		body := w.Body.String()
		if body != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})
}