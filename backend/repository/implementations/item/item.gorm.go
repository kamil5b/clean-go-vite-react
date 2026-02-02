package item

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMItemRepository is a GORM implementation of ItemRepository
type GORMItemRepository struct {
	db *gorm.DB
}

// ItemModel represents the items table schema
type ItemModel = entity.ItemEntity

// NewGORMItemRepository creates a new GORM item repository
func NewGORMItemRepository(db *gorm.DB) (*GORMItemRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&ItemModel{}); err != nil {
		return nil, err
	}

	return &GORMItemRepository{
		db: db,
	}, nil
}
