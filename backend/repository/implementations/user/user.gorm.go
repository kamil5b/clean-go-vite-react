package user

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMUserRepository is a GORM implementation of UserRepository
type GORMUserRepository struct {
	db *gorm.DB
}

// UserModel represents the users table schema
type UserModel = entity.UserEntity

// NewGORMUserRepository creates a new GORM user repository
func NewGORMUserRepository(db *gorm.DB) (*GORMUserRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil, err
	}

	return &GORMUserRepository{
		db: db,
	}, nil
}
