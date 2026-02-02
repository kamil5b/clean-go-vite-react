package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceItemEntity represents an invoice item in the system
type InvoiceItemEntity struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	InvoiceID  uuid.UUID `gorm:"index;default:0"`
	ItemID     uuid.UUID `gorm:"index;default:0"`
	Quantity   int       `gorm:"default:0"`
	UnitPrice  float64   `gorm:"column:unit_price;default:0"`
	TotalPrice float64   `gorm:"column:total_price;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Item       ItemEntity     `gorm:"foreignKey:ItemID"`
}

// TableName specifies the table name for InvoiceItemEntity
func (InvoiceItemEntity) TableName() string {
	return "invoice_items"
}
