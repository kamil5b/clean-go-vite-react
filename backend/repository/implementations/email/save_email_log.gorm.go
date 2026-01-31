package email

import (
	"context"
)

// SaveEmailLog saves an email log entry to GORM
func (r *GORMEmailRepository) SaveEmailLog(ctx context.Context, to, subject, body string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	emailLog := EmailLogModel{
		To:      to,
		Subject: subject,
		Body:    body,
		Status:  "sent",
	}

	if err := r.db.WithContext(ctx).Create(&emailLog).Error; err != nil {
		return err
	}

	return nil
}
