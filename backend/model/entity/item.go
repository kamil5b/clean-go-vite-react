package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ItemEntity represents an item in the system
type ItemEntity struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"default:''"`
	Desc      string    `gorm:"column:desc;default:''"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for ItemEntity
func (ItemEntity) TableName() string {
	return "items"
}
