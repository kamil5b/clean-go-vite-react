package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user entity.UserEntity) (*string, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.UserEntity, error)
	Update(ctx context.Context, id uuid.UUID, user entity.UserEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
}
