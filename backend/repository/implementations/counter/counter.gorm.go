package counter

import (
	"gorm.io/gorm"
)

// SQLiteCounterRepository is a SQLite implementation of CounterRepository
type SQLiteCounterRepository struct {
	db *gorm.DB
}

// CounterModel represents the counter table schema
type CounterModel struct {
	ID    uint `gorm:"primaryKey"`
	Value int
}

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
		db.Create(&CounterModel{ID: 1, Value: 0})
	}

	return &SQLiteCounterRepository{
		db: db,
	}, nil
}
