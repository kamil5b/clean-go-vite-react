package invoice

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindByID finds an invoice by ID with all related data
func (r *GORMInvoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.InvoiceEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var invoice entity.InvoiceEntity
	if err := r.db.WithContext(ctx).
		Preload("Items.Item").
		Preload("Tags").
		First(&invoice, id).Error; err != nil {
		return nil, err
	}

	return &invoice, nil
}
