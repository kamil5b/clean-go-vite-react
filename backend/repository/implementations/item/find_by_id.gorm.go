package item

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindByID finds an item by ID
func (r *GORMItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ItemEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var item entity.ItemEntity
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}
