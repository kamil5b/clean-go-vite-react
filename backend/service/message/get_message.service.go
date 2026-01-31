package message

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// GetMessage returns a greeting message using the response model
func (s *messageService) GetMessage(ctx context.Context) (*response.GetMessage, error) {
	message, err := s.repo.GetMessage(ctx, "default")
	if err != nil {
		return nil, err
	}
	return &response.GetMessage{
		Content: *message,
	}, nil
}
