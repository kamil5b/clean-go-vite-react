package message

import (
	"context"
	"fmt"
)

// GetMessage returns a stored message from SQLite
func (r *SQLiteMessageRepository) GetMessage(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	var message MessageModel
	if err := r.db.WithContext(ctx).Where("key = ?", "default").First(&message).Error; err != nil {
		return "", fmt.Errorf("message not found")
	}

	return message.Value, nil
}
