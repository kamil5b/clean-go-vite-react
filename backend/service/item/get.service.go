package item

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// GetByID gets an item by ID
func (s *itemService) GetByID(ctx context.Context, id uuid.UUID) (*response.ItemResponse, error) {
	item, err := s.itemRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.ItemResponse{
		ID:        item.ID,
		Name:      item.Name,
		Desc:      item.Desc,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

// GetAll gets all items with pagination
func (s *itemService) GetAll(ctx context.Context, page, limit int, search string) (*response.ItemPaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items, total, err := s.itemRepository.FindAll(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	itemResponses := make([]response.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = response.ItemResponse{
			ID:        item.ID,
			Name:      item.Name,
			Desc:      item.Desc,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &response.ItemPaginationResponse{
		Data: itemResponses,
		Meta: response.ItemPaginationMeta{
			TotalData: int(total),
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
		},
	}, nil
}
