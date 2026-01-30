package interfaces

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user map[string]interface{}) (string, error)
	FindByID(ctx context.Context, id string) (map[string]interface{}, error)
	Update(ctx context.Context, id string, user map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}
