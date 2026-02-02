package item

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Update updates an item by ID
func (r *GORMItemRepository) Update(ctx context.Context, id uuid.UUID, item entity.ItemEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Model(&entity.ItemEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name": item.Name,
			"desc": item.Desc,
		}).Error
}
