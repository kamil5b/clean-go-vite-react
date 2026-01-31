package message

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
)

// MessageService defines the interface for message operations
type MessageService interface {
	GetMessage(ctx context.Context) (*response.GetMessage, error)
}

// messageService is the concrete implementation of MessageService
type messageService struct {
	repo interfaces.MessageRepository
}

// NewMessageService creates a new instance of MessageService
func NewMessageService(repo interfaces.MessageRepository) MessageService {
	return &messageService{
		repo: repo,
	}
}
