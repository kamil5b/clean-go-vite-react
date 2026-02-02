package invoice

import (
	"context"

	"github.com/google/uuid"
)

// Delete deletes an invoice
func (s *invoiceService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if invoice exists
	_, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.invoiceRepository.Delete(ctx, id)
}
