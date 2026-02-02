package tag

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Create creates a new tag in GORM
func (r *GORMTagRepository) Create(ctx context.Context, tag entity.TagEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
		return nil, err
	}

	return &tag.ID, nil
}
