package message

import (
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// SQLiteMessageRepository is a SQLite implementation of MessageRepository
type SQLiteMessageRepository struct {
	db *gorm.DB
}

// MessageModel represents the message table schema
type MessageModel = entity.MessageEntity

// NewGORMMessageRepository creates a new SQLite message repository
func NewGORMMessageRepository(db *gorm.DB) (*SQLiteMessageRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&MessageModel{}); err != nil {
		return nil, err
	}

	// Initialize counter if it doesn't exist
	var count int64
	db.Model(&MessageModel{}).Count(&count)
	if count == 0 {
		db.Create(&MessageModel{ID: uuid.New(), Key: "default", Value: "Welcome to Clean Go Vite React!"})
	}

	return &SQLiteMessageRepository{
		db: db,
	}, nil
}
