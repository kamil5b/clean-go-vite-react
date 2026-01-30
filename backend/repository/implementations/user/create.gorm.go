package user

import (
	"context"
	"fmt"
)

// Create creates a new user in SQLite
func (r *SQLiteUserRepository) Create(ctx context.Context, user map[string]interface{}) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	userModel := UserModel{
		Name:  user["name"].(string),
		Email: user["email"].(string),
	}

	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("user_%d", userModel.ID), nil
}
