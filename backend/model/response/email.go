package response

import "github.com/google/uuid"

// Email represents an email record
type GetEmailLog struct {
	ID        uuid.UUID `json:"id"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Status    string    `json:"status"` // "pending", "sent", "failed"
	CreatedAt string    `json:"created_at"`
}
