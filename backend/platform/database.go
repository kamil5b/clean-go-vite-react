package platform

import (
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDatabase(cfg *Config) *gorm.DB {
	var dialector gorm.Dialector
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	if cfg.Database.Type == "postgres" {
		dialector = postgres.Open(cfg.Database.DSN)
	} else {
		dialector = sqlite.Open(cfg.Database.DSN)
	}
	db, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB() // Get the underlying generic *sql.DB
	if err != nil {
		panic("failed to get underlying sql.DB")
	}

	// Set connection pool settings from your config
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db
}
