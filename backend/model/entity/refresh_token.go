package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenEntity represents a refresh token stored in the database
type RefreshTokenEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
