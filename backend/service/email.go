package service

import (
	"context"
	"fmt"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

// emailService is the concrete implementation of EmailService
type emailService struct{}

// NewEmailService creates a new instance of EmailService
func NewEmailService() EmailService {
	return &emailService{}
}

// SendEmail sends an email (placeholder implementation)
func (s *emailService) SendEmail(ctx context.Context, to, subject, body string) error {
	// This is a placeholder implementation
	// In production, this would integrate with an email service (SendGrid, AWS SES, etc.)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// For now, just log the operation
	fmt.Printf("Sending email to %s with subject: %s\n", to, subject)
	return nil
}
