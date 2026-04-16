package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync/atomic"

    "github.com/gorilla/mux"
    "go-restapi-server/internal/models"
    "go-restapi-server/internal/store"
)

// PeopleHandler handles person-related operations
type PeopleHandler struct {
    store store.PersonStore
}

// NewPeopleHandler creates a new PeopleHandler
func NewPeopleHandler(store store.PersonStore) *PeopleHandler {
	return &PeopleHandler{store: store}
}

// GetPeopleEndpoint returns all people in the store
func (h *PeopleHandler) GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	// Parse query parameters
	pageStr := req.URL.Query().Get("page")
	limitStr := req.URL.Query().Get("limit")

	// Set defaults
	page := 1
	limit := 20
	maxLimit := 100

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid page parameter"})
			return
		}
	}

	// Parse limit parameter
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > maxLimit {
				limit = maxLimit
			} else {
				limit = l
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid limit parameter"})
			return
		}
	} else {
		// Default limit is 20
		limit = 20
	}

	// Get all people
	people := h.store.List()
	total := len(people)

	// Calculate pagination
	offset := (page - 1) * limit
	if offset >= total {
		page = (total + limit - 1) / limit
		if page < 1 {
			page = 1
		}
		offset = (page - 1) * limit
	}

	// Paginate the results
	var paginatedPeople []*models.Person
	if offset < total {
		end := offset + limit
		if end > total {
			end = total
		}
		paginatedPeople = people[offset:end]
	} else {
		// When the offset is beyond the total, return empty slice
		paginatedPeople = []*models.Person{}
	}

	// Calculate total pages
	pages := (total + limit - 1) / limit
	if pages < 1 {
		pages = 1
	}

	// Return paginated response
	response := map[string]interface{}{
		"data":  paginatedPeople,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": pages,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreatePersonEndpoint creates a new person
func (h *PeopleHandler) CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
    var person models.Person
    err := json.NewDecoder(req.Body).Decode(&person)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
        return
    }

	// Validate required fields
	if person.FirstName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "Firstname is required"})
		return
	}

	if person.LastName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"error": "Lastname is required"})
		return
	}

    // Generate unique ID using a simple, process-safe counter
    person.ID = nextID()

	// Create the person in the store
	err = h.store.Create(&person)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create person"})
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(person)
}

// nextID returns a monotonically increasing string ID.
var idCounter uint64

func nextID() string {
    n := atomic.AddUint64(&idCounter, 1)
    return strconv.FormatUint(n, 10)
}

// GetPersonEndpoint retrieves a person by ID
func (h *PeopleHandler) GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	person, err := h.store.Get(params["id"])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Person not found"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// UpdatePersonEndpoint updates a person by ID
func (h *PeopleHandler) UpdatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	// Get the person from the store
	person, err := h.store.Get(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Person not found"}`))
		return
	}

	// Decode the new person data
	var updatedPerson models.Person
	err = json.NewDecoder(req.Body).Decode(&updatedPerson)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Invalid JSON"}`))
		return
	}

	// Validate required fields
	if updatedPerson.FirstName == "" || updatedPerson.LastName == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error":"Firstname and Lastname are required"}`))
		return
	}

	// Update the person with new data
	person.FirstName = updatedPerson.FirstName
	person.LastName = updatedPerson.LastName
	person.Email = updatedPerson.Email
	person.Age = updatedPerson.Age

	// Save the updated person back to the store
	err = h.store.Update(person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to update person"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}

// DeletePersonEndpoint deletes a person by ID
func (h *PeopleHandler) DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	err := h.store.Delete(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Person not found"}`))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
