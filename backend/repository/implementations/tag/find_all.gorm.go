package tag

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindAll finds all tags with pagination and search
func (r *GORMTagRepository) FindAll(ctx context.Context, page, limit int, search string) ([]entity.TagEntity, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	var tags []entity.TagEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.TagEntity{})

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
	if err := query.Offset(offset).Limit(limit).Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}
