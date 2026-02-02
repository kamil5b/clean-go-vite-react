package invoice

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// GetByID gets an invoice by ID
func (s *invoiceService) GetByID(ctx context.Context, id uuid.UUID) (*response.InvoiceDetailResponse, error) {
	invoice, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(invoice), nil
}

// GetAll gets all invoices with pagination
func (s *invoiceService) GetAll(ctx context.Context, page, limit int, search string) (*response.InvoicePaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invoices, total, err := s.invoiceRepository.FindAll(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	invoiceList := make([]response.InvoiceListItem, len(invoices))
	for i, invoice := range invoices {
		tags := make([]response.TagResponse, len(invoice.Tags))
		for j, tag := range invoice.Tags {
			tags[j] = response.TagResponse{
				ID:       tag.ID,
				Name:     tag.Name,
				ColorHex: tag.ColorHex,
			}
		}

		invoiceList[i] = response.InvoiceListItem{
			ID:         invoice.ID,
			GrandPrice: invoice.GrandPrice,
			Tags:       tags,
			TotalItem:  len(invoice.Items),
			CreatedAt:  invoice.CreatedAt,
			UpdatedAt:  invoice.UpdatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &response.InvoicePaginationResponse{
		Data: invoiceList,
		Meta: response.InvoicePaginationMeta{
			TotalData: int(total),
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
		},
	}, nil
}
