package invoice

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Update updates an invoice by ID
func (r *GORMInvoiceRepository) Update(ctx context.Context, id uuid.UUID, invoice entity.InvoiceEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Model(&entity.InvoiceEntity{}).
		Where("id = ?", id).
		Update("grand_price", invoice.GrandPrice).Error
}
