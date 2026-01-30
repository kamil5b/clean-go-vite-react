package interfaces

import (
	"context"
)

// EmailRepository defines the interface for email data access
type EmailRepository interface {
	SaveEmailLog(ctx context.Context, to, subject, body string) error
	GetEmailLog(ctx context.Context, id string) (map[string]interface{}, error)
}
