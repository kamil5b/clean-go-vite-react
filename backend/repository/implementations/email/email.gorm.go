package email

import (
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// SQLiteEmailRepository is a SQLite implementation of EmailRepository
type SQLiteEmailRepository struct {
	db *gorm.DB
}

// EmailLogModel represents the email_logs table schema
type EmailLogModel entity.EmailLogEntity

// TableName specifies the table name for EmailLogModel
func (EmailLogModel) TableName() string {
	return "email_logs"
}

// NewSQLiteEmailRepository creates a new SQLite email repository
func NewSQLiteEmailRepository(db *gorm.DB) (*SQLiteEmailRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&EmailLogModel{}); err != nil {
		return nil, err
	}

	return &SQLiteEmailRepository{
		db: db,
	}, nil
}
