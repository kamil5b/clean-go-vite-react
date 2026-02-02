package invoice

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMInvoiceRepository is a GORM implementation of InvoiceRepository
type GORMInvoiceRepository struct {
	db *gorm.DB
}

// InvoiceModel represents the invoices table schema
type InvoiceModel = entity.InvoiceEntity

// NewGORMInvoiceRepository creates a new GORM invoice repository
func NewGORMInvoiceRepository(db *gorm.DB) (*GORMInvoiceRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&InvoiceModel{}, &entity.InvoiceItemEntity{}); err != nil {
		return nil, err
	}

	return &GORMInvoiceRepository{
		db: db,
	}, nil
}
