package item

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindAll finds all items with pagination and search
func (r *GORMItemRepository) FindAll(ctx context.Context, page, limit int, search string) ([]entity.ItemEntity, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	var items []entity.ItemEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.ItemEntity{})

	// Apply search filter
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
