package message

import (
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMMessageRepository is a GORM implementation of MessageRepository
type GORMMessageRepository struct {
	db *gorm.DB
}

// MessageModel represents the message table schema
type MessageModel = entity.MessageEntity

// NewGORMMessageRepository creates a new GORM message repository
func NewGORMMessageRepository(db *gorm.DB) (*GORMMessageRepository, error) {
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

	return &GORMMessageRepository{
		db: db,
	}, nil
}
