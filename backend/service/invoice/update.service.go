package invoice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Update updates an invoice
func (s *invoiceService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateInvoiceRequest) (*response.InvoiceDetailResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Check if invoice exists
	_, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete existing invoice items and tags
	if err := s.invoiceRepository.DeleteInvoiceItems(ctx, id); err != nil {
		return nil, err
	}
	if err := s.invoiceRepository.DeleteInvoiceTags(ctx, id); err != nil {
		return nil, err
	}

	// Update grand price
	invoice := entity.InvoiceEntity{
		GrandPrice: req.GrandPrice,
	}
	if err := s.invoiceRepository.Update(ctx, id, invoice); err != nil {
		return nil, err
	}

	// Create new invoice items with proper invoice ID
	invoiceItems := make([]entity.InvoiceItemEntity, len(req.Items))
	for i, item := range req.Items {
		totalPrice := float64(item.Quantity) * item.UnitPrice
		invoiceItems[i] = entity.InvoiceItemEntity{
			InvoiceID:  id,
			ItemID:     item.ItemID,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: totalPrice,
		}
	}

	// Re-create invoice with items and tags for associations
	invoiceWithRelations := entity.InvoiceEntity{
		ID:         id,
		GrandPrice: req.GrandPrice,
		Items:      invoiceItems,
	}

	// Build tags
	if len(req.Tags) > 0 {
		tags := make([]entity.TagEntity, len(req.Tags))
		for i, tagID := range req.Tags {
			tags[i] = entity.TagEntity{ID: tagID}
		}
		invoiceWithRelations.Tags = tags
	}

	// Save with associations - this creates items and tag associations
	s.invoiceRepository.Create(ctx, invoiceWithRelations)

	// Fetch updated invoice with all relations
	updated, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(updated), nil
}
