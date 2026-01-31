package user

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// SQLiteUserRepository is a SQLite implementation of UserRepository
type SQLiteUserRepository struct {
	db *gorm.DB
}

// UserModel represents the users table schema
type UserModel = entity.UserEntity

// NewGORMUserRepository creates a new SQLite user repository
func NewGORMUserRepository(db *gorm.DB) (*SQLiteUserRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{
		db: db,
	}, nil
}
