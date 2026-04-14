package store

import (
	"testing"

	"go-restapi-server/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryPersonStore_Create(t *testing.T) {
	store := NewInMemoryPersonStore()

	person := &models.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}

	err := store.Create(person)
	assert.NoError(t, err)

	// Verify the person was stored
	retrieved, err := store.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, person, retrieved)
}

func TestInMemoryPersonStore_Get(t *testing.T) {
	store := NewInMemoryPersonStore()

	// Test getting a non-existent person
	_, err := store.Get("non-existent")
	assert.Error(t, err)
	assert.IsType(t, &PersonNotFoundError{}, err)

	// Add a person
	person := &models.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	store.Create(person)

	// Test getting an existing person
	retrieved, err := store.Get("1")
	assert.NoError(t, err)
	assert.Equal(t, person, retrieved)
}

func TestInMemoryPersonStore_List(t *testing.T) {
	store := NewInMemoryPersonStore()

	// Test empty list
	persons := store.List()
	assert.Empty(t, persons)

	// Add some persons
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

	// Test list
	persons = store.List()
	assert.Len(t, persons, 2)

	// Verify persons are in the list (order might vary)
	found1 := false
	found2 := false
	for _, p := range persons {
		if p.ID == "1" {
			found1 = true
		}
		if p.ID == "2" {
			found2 = true
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}

func TestInMemoryPersonStore_Delete(t *testing.T) {
	store := NewInMemoryPersonStore()

	// Test deleting a non-existent person
	err := store.Delete("non-existent")
	assert.Error(t, err)
	assert.IsType(t, &PersonNotFoundError{}, err)

	// Add a person
	person := &models.Person{
		ID:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	store.Create(person)

	// Verify the person exists
	_, err = store.Get("1")
	assert.NoError(t, err)

	// Delete the person
	err = store.Delete("1")
	assert.NoError(t, err)

	// Verify the person was deleted
	_, err = store.Get("1")
	assert.Error(t, err)
	assert.IsType(t, &PersonNotFoundError{}, err)

	// Try to delete again - should fail
	err = store.Delete("1")
	assert.Error(t, err)
	assert.IsType(t, &PersonNotFoundError{}, err)
}
