package email

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// GORMEmailRepository is a GORM implementation of EmailRepository
type GORMEmailRepository struct {
	db *gorm.DB
}

// EmailLogModel represents the email_logs table schema
type EmailLogModel entity.EmailLogEntity

// TableName specifies the table name for EmailLogModel
func (EmailLogModel) TableName() string {
	return "email_logs"
}

// NewGORMEmailRepository creates a new GORM email repository
func NewGORMEmailRepository(db *gorm.DB) (*GORMEmailRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&EmailLogModel{}); err != nil {
		return nil, err
	}

	return &GORMEmailRepository{
		db: db,
	}, nil
}
