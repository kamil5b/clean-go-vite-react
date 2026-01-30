package repository

import (
	"context"
)

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	GetMessage(ctx context.Context) (string, error)
}

// EmailRepository defines the interface for email data access
type EmailRepository interface {
	SaveEmailLog(ctx context.Context, to, subject, body string) error
	GetEmailLog(ctx context.Context, id string) (map[string]interface{}, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user map[string]interface{}) (string, error)
	FindByID(ctx context.Context, id string) (map[string]interface{}, error)
	Update(ctx context.Context, id string, user map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

// CounterRepository defines the interface for counter data access
type CounterRepository interface {
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)
}
