package message

import (
	"context"
	"fmt"
)

// GetMessage returns a stored message from GORM
func (r *GORMMessageRepository) GetMessage(ctx context.Context, key string) (*string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var message MessageModel
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&message).Error; err != nil {
		return nil, fmt.Errorf("message not found")
	}
	resp := message.Value

	return &resp, nil
}
