package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
)

// EmailRepository defines the interface for email data access
type EmailRepository interface {
	SaveEmailLog(ctx context.Context, to, subject, body string) error
	GetEmailLog(ctx context.Context, id uuid.UUID) (*entity.EmailLogEntity, error)
}
