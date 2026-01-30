package service

import (
	"context"
	"testing"
	"time"
)

func TestEmailService_SendEmail_Success(t *testing.T) {
	// Arrange
	svc := NewEmailService()

	// Act
	err := svc.SendEmail(context.Background(), "test@example.com", "Test Subject", "Test Body")

	// Assert
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}
}

func TestEmailService_SendEmail_WithContext(t *testing.T) {
	// Arrange
	svc := NewEmailService()
	ctx := context.Background()

	// Act
	err := svc.SendEmail(ctx, "user@example.com", "Welcome", "Welcome to our service")

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEmailService_SendEmail_ContextCanceled(t *testing.T) {
	// Arrange
	svc := NewEmailService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := svc.SendEmail(ctx, "test@example.com", "Subject", "Body")

	// Assert
	if err == nil {
		t.Error("expected context error, got nil")
	}
}

func TestEmailService_SendEmail_ContextTimeout(t *testing.T) {
	// Arrange
	svc := NewEmailService()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Act
	err := svc.SendEmail(ctx, "test@example.com", "Subject", "Body")

	// Assert
	// Should succeed as service is fast
	if err != nil {
		t.Errorf("SendEmail failed: %v", err)
	}
}

func TestEmailService_SendEmail_EmptyFields(t *testing.T) {
	// Arrange
	svc := NewEmailService()

	tests := []struct {
		name    string
		to      string
		subject string
		body    string
	}{
		{"empty to", "", "Subject", "Body"},
		{"empty subject", "test@example.com", "", "Body"},
		{"empty body", "test@example.com", "Subject", ""},
		{"all empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			ctx := context.Background()
			err := svc.SendEmail(ctx, tt.to, tt.subject, tt.body)

			// Assert - Should still succeed (validation would be in handler/service wrapper)
			if err != nil {
				t.Errorf("SendEmail failed: %v", err)
			}
		})
	}
}
