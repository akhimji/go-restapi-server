package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"go-restapi-server/internal/models"
	"go-restapi-server/internal/store"
)

func TestGetPeopleEndpoint(t *testing.T) {
	// Create a test store with some data
	store := store.NewInMemoryPersonStore()

	// Add test people
	person1 := &models.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	person2 := &models.Person{
		ID:        "2",
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Age:       25,
	}
	store.Create(person1)
	store.Create(person2)

	// Create a people handler
	handler := NewPeopleHandler(store)

	// Test case: Get all people (default pagination)
	t.Run("get all people with default pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

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

		// Check response body is a JSON object with pagination data
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		// Check that the response has the expected structure
		if _, exists := response["data"]; !exists {
			t.Error("response missing 'data' field")
		}
		if _, exists := response["total"]; !exists {
			t.Error("response missing 'total' field")
		}
		if _, exists := response["page"]; !exists {
			t.Error("response missing 'page' field")
		}
		if _, exists := response["limit"]; !exists {
			t.Error("response missing 'limit' field")
		}
		if _, exists := response["pages"]; !exists {
			t.Error("response missing 'pages' field")
		}

		// Check that data is an array with 2 people
		data, ok := response["data"].([]interface{})
		if !ok {
			t.Error("data field is not an array")
		}
		if len(data) != 2 {
			t.Errorf("handler returned wrong number of people: got %v want %v",
				len(data), 2)
		}

		// Check that we got the right people
		person1Found := false
		person2Found := false
		for _, item := range data {
			if personMap, ok := item.(map[string]interface{}); ok {
				if personMap["id"] == "1" {
					person1Found = true
					if personMap["firstName"] != "John" || personMap["lastName"] != "Doe" || personMap["email"] != "john.doe@example.com" || personMap["age"] != 30.0 {
						t.Errorf("person 1 data mismatch")
					}
				}
				if personMap["id"] == "2" {
					person2Found = true
					if personMap["firstName"] != "Jane" || personMap["lastName"] != "Smith" || personMap["email"] != "jane.smith@example.com" || personMap["age"] != 25.0 {
						t.Errorf("person 2 data mismatch")
					}
				}
			}
		}
		if !person1Found || !person2Found {
			t.Errorf("not all people were returned")
		}

		// Check pagination metadata
		if total, ok := response["total"].(float64); ok {
			if int(total) != 2 {
				t.Errorf("total count is wrong: got %v want %v", int(total), 2)
			}
		}
		if page, ok := response["page"].(float64); ok {
			if int(page) != 1 {
				t.Errorf("page is wrong: got %v want %v", int(page), 1)
			}
		}
		if limit, ok := response["limit"].(float64); ok {
			if int(limit) != 20 {
				t.Errorf("limit is wrong: got %v want %v", int(limit), 20)
			}
		}
		if pages, ok := response["pages"].(float64); ok {
			if int(pages) != 1 {
				t.Errorf("pages count is wrong: got %v want %v", int(pages), 1)
			}
		}
	})

	// Test case: Get people with specific page and limit
	t.Run("get people with specific page and limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people?page=1&limit=1", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check response body
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		// Check that data contains only 1 person
		data, ok := response["data"].([]interface{})
		if !ok {
			t.Error("data field is not an array")
		}
		if len(data) != 1 {
			t.Errorf("handler returned wrong number of people: got %v want %v",
				len(data), 1)
		}

		// Check pagination metadata
		if total, ok := response["total"].(float64); ok {
			if int(total) != 2 {
				t.Errorf("total count is wrong: got %v want %v", int(total), 2)
			}
		}
		if page, ok := response["page"].(float64); ok {
			if int(page) != 1 {
				t.Errorf("page is wrong: got %v want %v", int(page), 1)
			}
		}
		if limit, ok := response["limit"].(float64); ok {
			if int(limit) != 1 {
				t.Errorf("limit is wrong: got %v want %v", int(limit), 1)
			}
		}
		if pages, ok := response["pages"].(float64); ok {
			if int(pages) != 2 {
				t.Errorf("pages count is wrong: got %v want %v", int(pages), 2)
			}
		}
	})

	// Test case: Get people with invalid page parameter
	t.Run("get people with invalid page parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people?page=0&limit=1", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case: Get people with invalid limit parameter
	t.Run("get people with invalid limit parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people?page=1&limit=0", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	// Test case: Get people with limit exceeding max
	t.Run("get people with limit exceeding max", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people?page=1&limit=150", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check response body
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

		// Check that limit was capped to 100
		if limit, ok := response["limit"].(float64); ok {
			if int(limit) != 100 {
				t.Errorf("limit should be capped to 100: got %v want %v", int(limit), 100)
			}
		}
	})
}

func TestGetPeopleEndpointEmptyStore(t *testing.T) {
	// Create an empty store
	store := store.NewInMemoryPersonStore()

	// Create a people handler
	handler := NewPeopleHandler(store)

	// Test case: Get all people from empty store
	t.Run("get all people from empty store", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/people", nil)
		w := httptest.NewRecorder()
		handler.GetPeopleEndpoint(w, req)

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

		// Check response body is a JSON object with pagination data
		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}

				
		// Check that data is an empty array
		data, exists := response["data"]
		if !exists {
			t.Error("data field missing")
		}
		if data == nil {
			t.Error("data field is nil")
		}
		dataSlice, ok := data.([]interface{})
		if !ok {
			t.Error("data field is not an array")
		}
		if len(dataSlice) != 0 {
			t.Errorf("handler returned wrong number of people: got %v want %v",
				len(dataSlice), 0)
		}

		// Check pagination metadata
		if total, ok := response["total"].(float64); ok {
			if int(total) != 0 {
				t.Errorf("total count is wrong: got %v want %v", int(total), 0)
			}
		}
		if page, ok := response["page"].(float64); ok {
			if int(page) != 1 {
				t.Errorf("page is wrong: got %v want %v", int(page), 1)
			}
		}
		if limit, ok := response["limit"].(float64); ok {
			if int(limit) != 20 {
				t.Errorf("limit is wrong: got %v want %v", int(limit), 20)
			}
		}
		if pages, ok := response["pages"].(float64); ok {
			if int(pages) != 1 {
				t.Errorf("pages count is wrong: got %v want %v", int(pages), 1)
			}
		}
	})
}

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

func TestCreatePersonEndpoint(t *testing.T) {
	// Create a test store
	store := store.NewInMemoryPersonStore()

	// Create a people handler
	handler := NewPeopleHandler(store)

	// Test case 1: Create valid person
	t.Run("create valid person", func(t *testing.T) {
		personData := map[string]interface{}{
			"firstName": "Jane",
			"lastName":  "Smith",
			"email":     "jane.smith@example.com",
			"age":       25,
		}

		jsonData, _ := json.Marshal(personData)
		req := httptest.NewRequest("POST", "/people", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.CreatePersonEndpoint(w, req)

		// Check status code
		if status := w.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
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

		// Verify fields are set correctly
		if responsePerson.FirstName != "Jane" {
			t.Errorf("handler returned wrong first name: got %v want %v",
				responsePerson.FirstName, "Jane")
		}
		if responsePerson.LastName != "Smith" {
			t.Errorf("handler returned wrong last name: got %v want %v",
				responsePerson.LastName, "Smith")
		}
		if responsePerson.Email != "jane.smith@example.com" {
			t.Errorf("handler returned wrong email: got %v want %v",
				responsePerson.Email, "jane.smith@example.com")
		}
		if responsePerson.Age != 25 {
			t.Errorf("handler returned wrong age: got %v want %v",
				responsePerson.Age, 25)
		}

		// Verify ID was generated
		if responsePerson.ID == "" {
			t.Errorf("handler did not generate ID")
		}
	})

	// Test case 2: Create person with invalid JSON
	t.Run("create person with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/people", bytes.NewBufferString(`{"firstName": "Jane", "lastName":}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.CreatePersonEndpoint(w, req)

		// Check status code
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		// Check error message
		expectedBody := `{"error":"Invalid JSON"}`
		body := w.Body.String()
		// Trim whitespace to handle possible newlines
		body = strings.TrimSpace(body)
		if body != expectedBody {
			t.Logf("Expected length: %d, Got length: %d", len(expectedBody), len(body))
			t.Logf("Expected: %q", expectedBody)
			t.Logf("Got: %q", body)
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})

	// Test case 3: Create person with missing firstName
	t.Run("create person with missing firstName", func(t *testing.T) {
		personData := map[string]interface{}{
			"lastName": "Smith",
			"email":    "jane.smith@example.com",
			"age":      25,
		}

		jsonData, _ := json.Marshal(personData)
		req := httptest.NewRequest("POST", "/people", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.CreatePersonEndpoint(w, req)

		// Check status code
		if status := w.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnprocessableEntity)
		}

		// Check error message
		expectedBody := `{"error":"Firstname is required"}`
		body := w.Body.String()
		// Trim whitespace to handle possible newlines
		body = strings.TrimSpace(body)
		if body != expectedBody {
			t.Logf("Expected: %s", expectedBody)
			t.Logf("Got: %s", body)
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})

	// Test case 4: Create person with missing lastName
	t.Run("create person with missing lastName", func(t *testing.T) {
		personData := map[string]interface{}{
			"firstName": "Jane",
			"email":     "jane.smith@example.com",
			"age":       25,
		}

		jsonData, _ := json.Marshal(personData)
		req := httptest.NewRequest("POST", "/people", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.CreatePersonEndpoint(w, req)

		// Check status code
		if status := w.Code; status != http.StatusUnprocessableEntity {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnprocessableEntity)
		}

		// Check error message
		expectedBody := `{"error":"Lastname is required"}`
		body := w.Body.String()
		// Trim whitespace to handle possible newlines
		body = strings.TrimSpace(body)
		if body != expectedBody {
			t.Logf("Expected: %s", expectedBody)
			t.Logf("Got: %s", body)
			t.Errorf("handler returned unexpected body: got %v want %v",
				body, expectedBody)
		}
	})
}

func TestDeletePersonEndpoint(t *testing.T) {
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

	// Test case 1: Delete existing person
	t.Run("delete existing person", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/people/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.DeletePersonEndpoint(w, req)

		if status := w.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNoContent)
		}

		// Verify person was actually deleted from store
		_, err := store.Get("1")
		if err == nil {
			t.Errorf("person was not deleted from store")
		}
	})

	// Test case 2: Delete non-existing person
	t.Run("delete non-existing person", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/people/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()
		handler.DeletePersonEndpoint(w, req)

		if status := w.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
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