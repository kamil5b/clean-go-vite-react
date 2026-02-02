package domain

import (
	"time"

	"github.com/google/uuid"
)

// Message represents a message in the system
type Message struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Item represents an item in the system
type Item struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
