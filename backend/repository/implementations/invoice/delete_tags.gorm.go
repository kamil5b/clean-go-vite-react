package invoice

import (
	"context"

	"github.com/google/uuid"
)

// DeleteInvoiceTags removes all tag associations for an invoice
func (r *GORMInvoiceRepository) DeleteInvoiceTags(ctx context.Context, invoiceID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.WithContext(ctx).Exec("DELETE FROM invoice_to_tags WHERE invoice_id = ?", invoiceID).Error
}
