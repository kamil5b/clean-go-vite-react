package model

// Message represents a message in the system
type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// User represents a user in the system
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Email represents an email record
type Email struct {
	ID        string `json:"id"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
	Status    string `json:"status"` // "pending", "sent", "failed"
	CreatedAt string `json:"created_at"`
	SentAt    string `json:"sent_at,omitempty"`
}

// Counter represents a counter in the system
type Counter struct {
	Value int `json:"value"`
}
