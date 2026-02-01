package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/domain"
)

// Logic defines the interface for domain business logic
type Logic interface {
	// Message
	GetMessage(ctx context.Context) (string, error)

	// Counter
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)

	// User
	RegisterUser(ctx context.Context, email, password, name string) (*domain.UserInfo, string, error)
	LoginUser(ctx context.Context, email, password string) (*domain.UserInfo, string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserInfo, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}
