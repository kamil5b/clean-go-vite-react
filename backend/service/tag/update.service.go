package tag

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Update updates a tag
func (s *tagService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateTagRequest) (*response.TagResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.ColorHex == "" {
		return nil, errors.New("color_hex is required")
	}

	// Check if tag exists
	_, err := s.tagRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	tag := entity.TagEntity{
		Name:     req.Name,
		ColorHex: req.ColorHex,
	}

	if err := s.tagRepository.Update(ctx, id, tag); err != nil {
		return nil, err
	}

	updated, err := s.tagRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.TagResponse{
		ID:        updated.ID,
		Name:      updated.Name,
		ColorHex:  updated.ColorHex,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}
