package tag

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMTagRepository is a GORM implementation of TagRepository
type GORMTagRepository struct {
	db *gorm.DB
}

// TagModel represents the tags table schema
type TagModel = entity.TagEntity

// NewGORMTagRepository creates a new GORM tag repository
func NewGORMTagRepository(db *gorm.DB) (*GORMTagRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&TagModel{}); err != nil {
		return nil, err
	}

	return &GORMTagRepository{
		db: db,
	}, nil
}
