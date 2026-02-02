package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceEntity represents an invoice in the system
type InvoiceEntity struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	GrandPrice float64   `gorm:"column:grand_price;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt      `gorm:"index"`
	Items      []InvoiceItemEntity `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`
	Tags       []TagEntity         `gorm:"many2many:invoice_to_tags;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for InvoiceEntity
func (InvoiceEntity) TableName() string {
	return "invoices"
}
