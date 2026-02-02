package item

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Update updates an item
func (s *itemService) Update(ctx context.Context, id uuid.UUID, req *request.UpdateItemRequest) (*response.ItemResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	// Check if item exists
	_, err := s.itemRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	item := entity.ItemEntity{
		Name: req.Name,
		Desc: req.Desc,
	}

	if err := s.itemRepository.Update(ctx, id, item); err != nil {
		return nil, err
	}

	updated, err := s.itemRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.ItemResponse{
		ID:        updated.ID,
		Name:      updated.Name,
		Desc:      updated.Desc,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
	}, nil
}
