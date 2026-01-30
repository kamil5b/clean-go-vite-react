package message

import "context"

// GetMessage returns a greeting message
func (s *messageService) GetMessage(ctx context.Context) (string, error) {
	return "Hello, from the golang World!", nil
}
