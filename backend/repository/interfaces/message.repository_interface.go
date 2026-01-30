package interfaces

import (
	"context"
)

// MessageRepository defines the interface for message data access
type MessageRepository interface {
	GetMessage(ctx context.Context) (string, error)
}
