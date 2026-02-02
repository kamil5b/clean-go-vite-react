package item

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
)

// ItemService defines the interface for item operations
type ItemService interface {
	Create(ctx context.Context, req *request.CreateItemRequest) (*response.ItemResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*response.ItemResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateItemRequest) (*response.ItemResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int, search string) (*response.ItemPaginationResponse, error)
}

// itemService is the concrete implementation of ItemService
type itemService struct {
	itemRepository interfaces.ItemRepository
}

// NewItemService creates a new instance of ItemService
func NewItemService(itemRepository interfaces.ItemRepository) ItemService {
	return &itemService{
		itemRepository: itemRepository,
	}
}
