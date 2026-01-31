package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// FindByID finds a user by ID in SQLite
func (r *SQLiteUserRepository) FindByID(ctx context.Context, id uuid.UUID) (user *entity.UserEntity, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return user, nil
}
