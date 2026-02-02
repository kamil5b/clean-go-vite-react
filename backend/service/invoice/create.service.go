package invoice

import (
	"context"
	"errors"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Create creates a new invoice
func (s *invoiceService) Create(ctx context.Context, req *request.CreateInvoiceRequest) (*response.InvoiceDetailResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Build invoice items
	invoiceItems := make([]entity.InvoiceItemEntity, len(req.Items))
	for i, item := range req.Items {
		totalPrice := float64(item.Quantity) * item.UnitPrice
		invoiceItems[i] = entity.InvoiceItemEntity{
			ItemID:     item.ItemID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: totalPrice,
		}
	}

	// Build tags
	tags := make([]entity.TagEntity, len(req.Tags))
	for i, tagID := range req.Tags {
		tags[i] = entity.TagEntity{ID: tagID}
	}

	invoice := entity.InvoiceEntity{
		GrandPrice: req.GrandPrice,
		Items:      invoiceItems,
		Tags:       tags,
	}

	id, err := s.invoiceRepository.Create(ctx, invoice)
	if err != nil {
		return nil, err
	}

	// Fetch created invoice with all relations
	created, err := s.invoiceRepository.FindByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(created), nil
}

func (s *invoiceService) toDetailResponse(invoice *entity.InvoiceEntity) *response.InvoiceDetailResponse {
	items := make([]response.InvoiceItemResponse, len(invoice.Items))
	for i, item := range invoice.Items {
		items[i] = response.InvoiceItemResponse{
			ID:     item.ID,
			ItemID: item.ItemID,
			Item: response.ItemResponse{
				ID:   item.Item.ID,
				Name: item.Item.Name,
				Desc: item.Item.Desc,
			},
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: item.TotalPrice,
		}
	}

	tags := make([]response.TagResponse, len(invoice.Tags))
	for i, tag := range invoice.Tags {
		tags[i] = response.TagResponse{
			ID:       tag.ID,
			Name:     tag.Name,
			ColorHex: tag.ColorHex,
		}
	}

	return &response.InvoiceDetailResponse{
		ID:         invoice.ID,
		GrandPrice: invoice.GrandPrice,
		Items:      items,
		Tags:       tags,
		CreatedAt:  invoice.CreatedAt,
		UpdatedAt:  invoice.UpdatedAt,
	}
}
