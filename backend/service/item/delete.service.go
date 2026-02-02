package item

import (
	"context"

	"github.com/google/uuid"
)

// Delete deletes an item
func (s *itemService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if item exists
	_, err := s.itemRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.itemRepository.Delete(ctx, id)
}
