package invoice

import (
	"context"
	"strconv"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindAll finds all invoices with pagination and search
func (r *GORMInvoiceRepository) FindAll(ctx context.Context, page, limit int, search string) ([]entity.InvoiceEntity, int64, error) {
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	var invoices []entity.InvoiceEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.InvoiceEntity{}).
		Preload("Tags").
		Preload("Items")

	// Apply search filter (search by ID)
	if search != "" {
		if id, err := strconv.ParseUint(search, 10, 32); err == nil {
			query = query.Where("id = ?", id)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&invoices).Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}
