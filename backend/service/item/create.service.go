package item

import (
	"context"
	"errors"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Create creates a new item
func (s *itemService) Create(ctx context.Context, req *request.CreateItemRequest) (*response.ItemResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	item := entity.ItemEntity{
		Name: req.Name,
		Desc: req.Desc,
	}

	id, err := s.itemRepository.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	created, err := s.itemRepository.FindByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return &response.ItemResponse{
		ID:        created.ID,
		Name:      created.Name,
		Desc:      created.Desc,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}
