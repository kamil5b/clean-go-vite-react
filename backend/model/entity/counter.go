package entity

import (
	"time"

	"github.com/google/uuid"
)

// CounterModel represents the counter table schema
type CounterEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Value     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
