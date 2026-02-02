package invoice

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
)

// InvoiceService defines the interface for invoice operations
type InvoiceService interface {
	Create(ctx context.Context, req *request.CreateInvoiceRequest) (*response.InvoiceDetailResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*response.InvoiceDetailResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateInvoiceRequest) (*response.InvoiceDetailResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int, search string) (*response.InvoicePaginationResponse, error)
}

// invoiceService is the concrete implementation of InvoiceService
type invoiceService struct {
	invoiceRepository interfaces.InvoiceRepository
	tagRepository     interfaces.TagRepository
}

// NewInvoiceService creates a new instance of InvoiceService
func NewInvoiceService(invoiceRepository interfaces.InvoiceRepository, tagRepository interfaces.TagRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository: invoiceRepository,
		tagRepository:     tagRepository,
	}
}
