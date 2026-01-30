package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/kamil5b/clean-go-vite-react/backend/service/email"
	"github.com/kamil5b/clean-go-vite-react/backend/task"
)

// EmailProcessor handles email notification tasks
type EmailProcessor struct {
	service email.EmailService
}

// NewEmailProcessor creates a new email task processor
func NewEmailProcessor(svc email.EmailService) *EmailProcessor {
	return &EmailProcessor{
		service: svc,
	}
}

// ProcessTask processes an email notification task
func (p *EmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload task.EmailNotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	return p.service.SendEmail(ctx, payload.To, payload.Subject, payload.Body)
}
