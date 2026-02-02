package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// ItemRepository defines the interface for item data access
type ItemRepository interface {
	Create(ctx context.Context, item entity.ItemEntity) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ItemEntity, error)
	Update(ctx context.Context, id uuid.UUID, item entity.ItemEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int, search string) ([]entity.ItemEntity, int64, error)
}
