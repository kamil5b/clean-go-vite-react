package user

import (
	"context"

	"github.com/google/uuid"
)

// Delete deletes a user in GORM
func (r *GORMUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := r.db.WithContext(ctx).Delete(&UserModel{}, id).Error; err != nil {
		return err
	}

	return nil
}
