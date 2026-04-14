package store

import (
	"sync"

	"go-restapi-server/internal/models"
)

// PersonStore interface defines the operations for person data storage
type PersonStore interface {
	Create(person *models.Person) error
	Get(id string) (*models.Person, error)
	List() []*models.Person
	Delete(id string) error
	Update(person *models.Person) error
}

// InMemoryPersonStore is an in-memory implementation of PersonStore
type InMemoryPersonStore struct {
	data map[string]*models.Person
	mu   sync.RWMutex
}

// NewInMemoryPersonStore creates and returns a new InMemoryPersonStore
func NewInMemoryPersonStore() *InMemoryPersonStore {
	return &InMemoryPersonStore{
		data: make(map[string]*models.Person),
	}
}

// Create adds a new person to the store
func (s *InMemoryPersonStore) Create(person *models.Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[person.ID] = person
	return nil
}

// Update updates an existing person in the store
func (s *InMemoryPersonStore) Update(person *models.Person) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[person.ID]
	if !exists {
		return &PersonNotFoundError{ID: person.ID}
	}

	s.data[person.ID] = person
	return nil
}

// Get retrieves a person by ID
func (s *InMemoryPersonStore) Get(id string) (*models.Person, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	person, exists := s.data[id]
	if !exists {
		return nil, &PersonNotFoundError{ID: id}
	}

	return person, nil
}

// List returns all persons in the store
func (s *InMemoryPersonStore) List() []*models.Person {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of the slice to avoid external mutation
	persons := make([]*models.Person, 0, len(s.data))
	for _, person := range s.data {
		persons = append(persons, person)
	}

	return persons
}

// Delete removes a person by ID
func (s *InMemoryPersonStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[id]
	if !exists {
		return &PersonNotFoundError{ID: id}
	}

	delete(s.data, id)
	return nil
}

// PersonNotFoundError is returned when a person is not found
type PersonNotFoundError struct {
	ID string
}

func (e *PersonNotFoundError) Error() string {
	return "person not found: " + e.ID
}
