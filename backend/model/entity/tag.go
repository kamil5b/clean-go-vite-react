package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TagEntity represents a tag in the system
type TagEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"default:''"`
	ColorHex  string    `gorm:"column:color_hex;default:'#000000'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for TagEntity
func (TagEntity) TableName() string {
	return "tags"
}
