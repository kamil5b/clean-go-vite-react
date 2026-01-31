package entity

import (
	"time"

	"github.com/google/uuid"
)

// MessageEntity represents a message in the system
type MessageEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
