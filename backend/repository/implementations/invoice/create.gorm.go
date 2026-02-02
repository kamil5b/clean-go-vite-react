package invoice

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Create creates a new invoice in GORM
func (r *GORMInvoiceRepository) Create(ctx context.Context, invoice entity.InvoiceEntity) (*uuid.UUID, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Create(&invoice).Error; err != nil {
		return nil, err
	}

	return &invoice.ID, nil
}
