package invoice

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// DeleteInvoiceItems deletes all items for an invoice
func (r *GORMInvoiceRepository) DeleteInvoiceItems(ctx context.Context, invoiceID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).
		Where("invoice_id = ?", invoiceID).
		Delete(&entity.InvoiceItemEntity{}).Error
}
