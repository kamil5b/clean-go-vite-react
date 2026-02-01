package infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the database connection and provides data access methods
type DB struct {
	conn *gorm.DB
}

// Config holds database configuration
type Config struct {
	Type            string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewDB creates a new database connection and runs migrations
func NewDB(cfg Config) (*DB, error) {
	var dialector gorm.Dialector
	dbConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.Type == "postgres" {
		dialector = postgres.Open(cfg.DSN)
	} else {
		dialector = sqlite.Open(cfg.DSN)
	}

	conn, err := gorm.Open(dialector, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	db := &DB{conn: conn}

	// Auto-migrate schemas
	if err := db.migrate(); err != nil {
		return nil, err
	}

	// Initialize default data
	if err := db.seed(); err != nil {
		return nil, err
	}

	return db, nil
}

// migrate runs database migrations
func (db *DB) migrate() error {
	return db.conn.AutoMigrate(
		&domain.Message{},
		&domain.Counter{},
		&domain.User{},
	)
}

// seed initializes default data
func (db *DB) seed() error {
	// Seed default message
	var messageCount int64
	db.conn.Model(&domain.Message{}).Count(&messageCount)
	if messageCount == 0 {
		msg := domain.Message{
			ID:    uuid.New(),
			Key:   "default",
			Value: "Welcome to Clean Go Vite React!",
		}
		if err := db.conn.Create(&msg).Error; err != nil {
			return err
		}
	}

	// Seed default counter
	var counterCount int64
	db.conn.Model(&domain.Counter{}).Count(&counterCount)
	if counterCount == 0 {
		counter := domain.Counter{
			ID:    uuid.New(),
			Value: 0,
		}
		if err := db.conn.Create(&counter).Error; err != nil {
			return err
		}
	}

	return nil
}

// Message operations

// GetMessage retrieves a message by key
func (db *DB) GetMessage(ctx context.Context, key string) (*domain.Message, error) {
	var message domain.Message
	if err := db.conn.WithContext(ctx).Where("key = ?", key).First(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, err
	}
	return &message, nil
}

// Counter operations

// GetCounter retrieves the current counter value
func (db *DB) GetCounter(ctx context.Context) (int, error) {
	var counter domain.Counter
	if err := db.conn.WithContext(ctx).First(&counter).Error; err != nil {
		return 0, err
	}
	return counter.Value, nil
}

// IncrementCounter atomically increments the counter and returns the new value
func (db *DB) IncrementCounter(ctx context.Context) (int, error) {
	var counter domain.Counter

	err := db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&counter).Error; err != nil {
			return err
		}
		counter.Value++
		return tx.Save(&counter).Error
	})

	if err != nil {
		return 0, err
	}

	return counter.Value, nil
}

// User operations

// CreateUser creates a new user
func (db *DB) CreateUser(ctx context.Context, user *domain.User) error {
	return db.conn.WithContext(ctx).Create(user).Error
}

// FindUserByEmail finds a user by email
func (db *DB) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := db.conn.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if not found
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByID finds a user by ID
func (db *DB) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := db.conn.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
