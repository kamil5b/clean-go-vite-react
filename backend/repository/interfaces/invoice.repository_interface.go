package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
	Create(ctx context.Context, invoice entity.InvoiceEntity) (*uuid.UUID, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.InvoiceEntity, error)
	Update(ctx context.Context, id uuid.UUID, invoice entity.InvoiceEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, page, limit int, search string) ([]entity.InvoiceEntity, int64, error)
	DeleteInvoiceItems(ctx context.Context, invoiceID uuid.UUID) error
	DeleteInvoiceTags(ctx context.Context, invoiceID uuid.UUID) error
}
