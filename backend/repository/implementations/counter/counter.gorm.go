package counter

import (
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"gorm.io/gorm"
)

// SQLiteCounterRepository is a SQLite implementation of CounterRepository
type SQLiteCounterRepository struct {
	db *gorm.DB
}

type CounterModel entity.CounterEntity

// TableName specifies the table name for CounterModel
func (CounterModel) TableName() string {
	return "counters"
}

// NewSQLiteCounterRepository creates a new SQLite counter repository
func NewSQLiteCounterRepository(db *gorm.DB) (*SQLiteCounterRepository, error) {
	// Auto-migrate the schema
	if err := db.AutoMigrate(&CounterModel{}); err != nil {
		return nil, err
	}

	// Initialize counter if it doesn't exist
	var count int64
	db.Model(&CounterModel{}).Count(&count)
	if count == 0 {
		db.Create(&CounterModel{ID: uuid.New(), Value: 0})
	}

	return &SQLiteCounterRepository{
		db: db,
	}, nil
}
