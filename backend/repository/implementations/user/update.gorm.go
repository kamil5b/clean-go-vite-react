package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// Update updates a user in GORM
func (r *GORMUserRepository) Update(ctx context.Context, id uuid.UUID, user entity.UserEntity) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).
		Model(&UserModel{}).Where("id = ?", id).
		Updates(&user).
		Error; err != nil {
		return err
	}

	return nil
}
