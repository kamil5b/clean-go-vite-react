package email

import (
	"context"
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
