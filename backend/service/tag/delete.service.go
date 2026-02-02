package tag

import (
	"context"

	"github.com/google/uuid"
)

// Delete deletes a tag
func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if tag exists
	_, err := s.tagRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.tagRepository.Delete(ctx, id)
}
