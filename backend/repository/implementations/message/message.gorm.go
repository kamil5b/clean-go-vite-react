package message

import (
	"gorm.io/gorm"
)

// SQLiteMessageRepository is a SQLite implementation of MessageRepository
type SQLiteMessageRepository struct {
	db *gorm.DB
}

// MessageModel represents the message table schema
type MessageModel struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex"`
	Value string
}

// TableName specifies the table name for MessageModel
func (MessageModel) TableName() string {
	return "messages"
}

// NewSQLiteMessageRepository creates a new SQLite message repository
func NewSQLiteMessageRepository(db *gorm.DB) (*SQLiteMessageRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&MessageModel{}); err != nil {
		return nil, err
	}

	return &SQLiteMessageRepository{
		db: db,
	}, nil
}
