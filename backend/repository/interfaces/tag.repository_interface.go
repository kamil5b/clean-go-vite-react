package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// TagRepository defines the interface for tag data access
type TagRepository interface {
	Create(ctx context.Context, tag entity.TagEntity) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.TagEntity, error)
	Update(ctx context.Context, id uuid.UUID, tag entity.TagEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int, search string) ([]entity.TagEntity, int64, error)
}
