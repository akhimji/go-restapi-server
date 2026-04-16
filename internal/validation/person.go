package validation

import (
	"fmt"
	"regexp"
	"strings"

	"go-restapi-server/internal/models"
)

// PersonValidationError represents a validation error for a person
type PersonValidationError struct {
	Field   string
	Message string
}

func (e *PersonValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// PersonValidationErrors represents multiple validation errors
type PersonValidationErrors struct {
	Errors []PersonValidationError
}

func (e *PersonValidationErrors) Error() string {
	var msgs []string
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return "validation failed: " + strings.Join(msgs, "; ")
}

// ValidatePerson validates a person object according to business rules
func ValidatePerson(person *models.Person) error {
	var errors []PersonValidationError

	// Validate name (required, 1-100 characters)
	if person.FirstName == "" {
		errors = append(errors, PersonValidationError{
			Field:   "firstName",
			Message: "is required",
		})
	} else if len(person.FirstName) < 1 || len(person.FirstName) > 100 {
		errors = append(errors, PersonValidationError{
			Field:   "firstName",
			Message: "must be 1-100 characters",
		})
	}

	if person.LastName == "" {
		errors = append(errors, PersonValidationError{
			Field:   "lastName",
			Message: "is required",
		})
	} else if len(person.LastName) < 1 || len(person.LastName) > 100 {
		errors = append(errors, PersonValidationError{
			Field:   "lastName",
			Message: "must be 1-100 characters",
		})
	}

	// Validate age (if present, must be 0-150)
	if person.Age < 0 || person.Age > 150 {
		errors = append(errors, PersonValidationError{
			Field:   "age",
			Message: "must be between 0 and 150",
		})
	}

	// Validate email (if present, must match basic email format)
	if person.Email != "" {
		// Basic email validation using regex
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(person.Email) {
			errors = append(errors, PersonValidationError{
				Field:   "email",
				Message: "must be a valid email address",
			})
		}
	}

	if len(errors) > 0 {
		return &PersonValidationErrors{Errors: errors}
	}

	return nil
}