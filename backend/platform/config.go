package platform

import (
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Gorm            *gorm.DB
	Type            string
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	DB       int
	Password string
}

// NewConfig loads configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			Host:         getEnv("SERVER_HOST", ""),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "dev.db"),
			Type:            getEnv("DATABASE_TYPE", "sqlite"),
			MaxOpenConns:    getEnvInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			DB:       getEnvInt("REDIS_DB", 0),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
	}
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
