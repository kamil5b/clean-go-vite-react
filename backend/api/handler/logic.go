package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/domain"
)

// Logic defines the interface for domain business logic
type Logic interface {
	// Message
	GetMessage(ctx context.Context) (string, error)

	// User
	RegisterUser(ctx context.Context, email, password, name string) (*domain.UserInfo, string, error)
	LoginUser(ctx context.Context, email, password string) (*domain.UserInfo, string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserInfo, error)
	ValidateToken(tokenString string) (uuid.UUID, error)

	// Item
	CreateItem(ctx context.Context, title, description string, userID uuid.UUID) (*domain.Item, error)
	GetItems(ctx context.Context, userID uuid.UUID) ([]*domain.Item, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*domain.Item, error)
	UpdateItem(ctx context.Context, id uuid.UUID, title, description string) (*domain.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
}

// ItemToInfo is a helper to convert Item to ItemInfo
func ItemToInfo(item *domain.Item) *domain.ItemInfo {
	return domain.ItemToInfo(item)
}

// ItemsToInfoList is a helper to convert Items slice to ItemInfo slice
func ItemsToInfoList(items []*domain.Item) []*domain.ItemInfo {
	return domain.ItemsToInfoList(items)
}
