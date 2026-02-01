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

// Counter represents a counter in the system
type Counter struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Value     int
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
