package models

// Person represents a person entity
type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
}

// PatchPerson represents a partial update to a person entity
type PatchPerson struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Email     *string `json:"email,omitempty"`
	Age       *int    `json:"age,omitempty"`
}
