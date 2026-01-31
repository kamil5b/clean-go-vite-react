package user

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// FindByEmail finds a user by email in GORM
func (r *GORMUserRepository) FindByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var user entity.UserEntity
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
