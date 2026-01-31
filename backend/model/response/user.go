package response

import "github.com/google/uuid"

// User represents a user in the system
type GetUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
