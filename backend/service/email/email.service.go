package email

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(ctx context.Context, req *request.SaveEmailRequest) (*response.GetEmailLog, error)
}

// emailService is the concrete implementation of EmailService
type emailService struct{}

// NewEmailService creates a new instance of EmailService
func NewEmailService() EmailService {
	return &emailService{}
}
