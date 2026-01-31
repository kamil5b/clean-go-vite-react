package email

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// SendEmail sends an email using the request model and returns a response
func (s *emailService) SendEmail(ctx context.Context, req *request.SaveEmailRequest) (*response.GetEmailLog, error) {
	// Validate context
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// For now, just log the operation
	fmt.Printf("Sending email to %s with subject: %s\n", req.To, req.Subject)

	// Build response
	emailLog := &response.GetEmailLog{
		ID:        uuid.New(),
		To:        req.To,
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    "sent",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return emailLog, nil
}
