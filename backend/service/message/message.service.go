package message

import (
	"context"
)

// MessageService defines the interface for message operations
type MessageService interface {
	GetMessage(ctx context.Context) (string, error)
}

// messageService is the concrete implementation of MessageService
type messageService struct{}

// NewMessageService creates a new instance of MessageService
func NewMessageService() MessageService {
	return &messageService{}
}
