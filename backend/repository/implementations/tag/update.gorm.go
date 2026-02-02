package tag

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Update updates a tag by ID
func (r *GORMTagRepository) Update(ctx context.Context, id uuid.UUID, tag entity.TagEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Model(&entity.TagEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":      tag.Name,
			"color_hex": tag.ColorHex,
		}).Error
}
