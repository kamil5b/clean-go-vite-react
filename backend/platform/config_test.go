package platform

import (
	"os"
	"testing"
	"time"
)

func TestNewConfig_DefaultValues(t *testing.T) {
	// Clear environment variables
	clearEnv()

	// Act
	cfg := NewConfig()

	// Assert
	if cfg == nil {
		t.Fatal("expected config to be created")
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Redis.Host != "localhost" {
		t.Errorf("expected default Redis host 'localhost', got %q", cfg.Redis.Host)
	}

	if cfg.Redis.Port != 6379 {
		t.Errorf("expected default Redis port 6379, got %d", cfg.Redis.Port)
	}

	if !cfg.Asynq.Enabled {
		t.Error("expected Asynq to be enabled by default")
	}
}

func TestNewConfig_EnvironmentOverrides(t *testing.T) {
	// Clear and set environment variables
	clearEnv()
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("REDIS_HOST", "redis.example.com")
	os.Setenv("REDIS_PORT", "6380")
	defer clearEnv()

	// Act
	cfg := NewConfig()

	// Assert
	if cfg.Server.Port != 9000 {
		t.Errorf("expected port 9000, got %d", cfg.Server.Port)
	}

	if cfg.Redis.Host != "redis.example.com" {
		t.Errorf("expected Redis host 'redis.example.com', got %q", cfg.Redis.Host)
	}

	if cfg.Redis.Port != 6380 {
		t.Errorf("expected Redis port 6380, got %d", cfg.Redis.Port)
	}
}

func TestNewConfig_DatabaseConfig(t *testing.T) {
	clearEnv()
	os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")
	os.Setenv("DATABASE_MAX_OPEN_CONNS", "50")
	os.Setenv("DATABASE_MAX_IDLE_CONNS", "10")
	defer clearEnv()

	// Act
	cfg := NewConfig()

	// Assert
	if cfg.Database.DSN != "postgres://user:pass@localhost/db" {
		t.Errorf("expected DSN, got %q", cfg.Database.DSN)
	}

	if cfg.Database.MaxOpenConns != 50 {
		t.Errorf("expected max open conns 50, got %d", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.MaxIdleConns != 10 {
		t.Errorf("expected max idle conns 10, got %d", cfg.Database.MaxIdleConns)
	}
}

func TestNewConfig_AsynqConfig(t *testing.T) {
	clearEnv()
	os.Setenv("ASYNQ_ENABLED", "true")
	os.Setenv("ASYNQ_REDIS_ADDR", "asynq.redis:6379")
	os.Setenv("ASYNQ_CONCURRENCY", "20")
	os.Setenv("ASYNQ_MAX_RETRIES", "5")
	defer clearEnv()

	// Act
	cfg := NewConfig()

	// Assert
	if !cfg.Asynq.Enabled {
		t.Error("expected Asynq to be enabled")
	}

	if cfg.Asynq.RedisAddr != "asynq.redis:6379" {
		t.Errorf("expected Redis addr 'asynq.redis:6379', got %q", cfg.Asynq.RedisAddr)
	}

	if cfg.Asynq.Concurrency != 20 {
		t.Errorf("expected concurrency 20, got %d", cfg.Asynq.Concurrency)
	}

	if cfg.Asynq.MaxRetries != 5 {
		t.Errorf("expected max retries 5, got %d", cfg.Asynq.MaxRetries)
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := getEnv("TEST_ENV_VAR", "default")

	if result != "test_value" {
		t.Errorf("expected 'test_value', got %q", result)
	}
}

func TestGetEnv_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_VAR")

	result := getEnv("NONEXISTENT_VAR", "default_value")

	if result != "default_value" {
		t.Errorf("expected 'default_value', got %q", result)
	}
}

func TestGetEnvInt_WithValue(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := getEnvInt("TEST_INT", 0)

	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestGetEnvInt_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_INT")

	result := getEnvInt("NONEXISTENT_INT", 99)

	if result != 99 {
		t.Errorf("expected 99, got %d", result)
	}
}

func TestGetEnvInt_InvalidValue(t *testing.T) {
	os.Setenv("INVALID_INT", "not_a_number")
	defer os.Unsetenv("INVALID_INT")

	result := getEnvInt("INVALID_INT", 77)

	if result != 77 {
		t.Errorf("expected default 77, got %d", result)
	}
}

func TestGetEnvBool_True(t *testing.T) {
	tests := []struct {
		value string
	}{
		{"true"},
		{"1"},
		{"yes"},
	}

	for _, tt := range tests {
		os.Setenv("TEST_BOOL", tt.value)
		result := getEnvBool("TEST_BOOL", false)
		if !result {
			t.Errorf("expected true for value %q", tt.value)
		}
		os.Unsetenv("TEST_BOOL")
	}
}

func TestGetEnvBool_False(t *testing.T) {
	tests := []struct {
		value string
	}{
		{"false"},
		{"0"},
		{"no"},
	}

	for _, tt := range tests {
		os.Setenv("TEST_BOOL", tt.value)
		result := getEnvBool("TEST_BOOL", true)
		if result {
			t.Errorf("expected false for value %q", tt.value)
		}
		os.Unsetenv("TEST_BOOL")
	}
}

func TestGetEnvBool_EmptyString(t *testing.T) {
	os.Setenv("TEST_BOOL", "")
	defer os.Unsetenv("TEST_BOOL")

	result := getEnvBool("TEST_BOOL", true)

	// Empty string returns default
	if !result {
		t.Error("expected default true for empty string")
	}
}

func TestGetEnvBool_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_BOOL")

	result := getEnvBool("NONEXISTENT_BOOL", true)

	if !result {
		t.Error("expected default true")
	}
}

func TestGetEnvDuration_WithValue(t *testing.T) {
	os.Setenv("TEST_DURATION", "30s")
	defer os.Unsetenv("TEST_DURATION")

	result := getEnvDuration("TEST_DURATION", 0)

	if result != 30*time.Second {
		t.Errorf("expected 30s, got %v", result)
	}
}

func TestGetEnvDuration_WithDefault(t *testing.T) {
	os.Unsetenv("NONEXISTENT_DURATION")

	result := getEnvDuration("NONEXISTENT_DURATION", 5*time.Minute)

	if result != 5*time.Minute {
		t.Errorf("expected 5m, got %v", result)
	}
}

func TestGetEnvDuration_InvalidValue(t *testing.T) {
	os.Setenv("INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("INVALID_DURATION")

	result := getEnvDuration("INVALID_DURATION", 10*time.Second)

	if result != 10*time.Second {
		t.Errorf("expected default 10s, got %v", result)
	}
}

func TestNewConfig_AllServerSettings(t *testing.T) {
	clearEnv()
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_READ_TIMEOUT", "20s")
	os.Setenv("SERVER_WRITE_TIMEOUT", "25s")
	os.Setenv("SERVER_IDLE_TIMEOUT", "90s")
	defer clearEnv()

	cfg := NewConfig()

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected host '0.0.0.0', got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("expected port 3000, got %d", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 20*time.Second {
		t.Errorf("expected read timeout 20s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 25*time.Second {
		t.Errorf("expected write timeout 25s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 90*time.Second {
		t.Errorf("expected idle timeout 90s, got %v", cfg.Server.IdleTimeout)
	}
}

func TestNewConfig_RedisAuth(t *testing.T) {
	clearEnv()
	os.Setenv("REDIS_HOST", "redis.prod")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_DB", "2")
	os.Setenv("REDIS_PASSWORD", "secret123")
	defer clearEnv()

	cfg := NewConfig()

	if cfg.Redis.Host != "redis.prod" {
		t.Errorf("expected host 'redis.prod'")
	}
	if cfg.Redis.DB != 2 {
		t.Errorf("expected db 2, got %d", cfg.Redis.DB)
	}
	if cfg.Redis.Password != "secret123" {
		t.Errorf("expected password 'secret123'")
	}
}

func TestNewConfig_AsynqDisabled(t *testing.T) {
	clearEnv()
	os.Setenv("ASYNQ_ENABLED", "false")
	defer clearEnv()

	cfg := NewConfig()

	if cfg.Asynq.Enabled {
		t.Error("expected Asynq to be disabled")
	}
}

// Helper function to clear all relevant environment variables
func clearEnv() {
	vars := []string{
		"SERVER_PORT", "SERVER_HOST", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"DATABASE_DSN", "DATABASE_MAX_OPEN_CONNS", "DATABASE_MAX_IDLE_CONNS", "DATABASE_CONN_MAX_LIFETIME",
		"REDIS_HOST", "REDIS_PORT", "REDIS_DB", "REDIS_PASSWORD",
		"ASYNQ_ENABLED", "ASYNQ_REDIS_ADDR", "ASYNQ_CONCURRENCY", "ASYNQ_MAX_RETRIES",
		"ASYNQ_RETRY_DELAY_MIN", "ASYNQ_RETRY_DELAY_MAX",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
}
