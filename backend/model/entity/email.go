package entity

import (
	"time"

	"github.com/google/uuid"
)

// EmailLogEntity represents the email_logs table schema
type EmailLogEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	To        string
	Subject   string
	Body      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
