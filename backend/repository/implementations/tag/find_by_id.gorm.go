package tag

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindByID finds a tag by ID
func (r *GORMTagRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.TagEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var tag entity.TagEntity
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}
