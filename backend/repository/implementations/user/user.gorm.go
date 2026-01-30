package user

import (
	"gorm.io/gorm"
)

// SQLiteUserRepository is a SQLite implementation of UserRepository
type SQLiteUserRepository struct {
	db *gorm.DB
}

// UserModel represents the users table schema
type UserModel struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// NewSQLiteUserRepository creates a new SQLite user repository
func NewSQLiteUserRepository(db *gorm.DB) (*SQLiteUserRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil, err
	}

	return &SQLiteUserRepository{
		db: db,
	}, nil
}
