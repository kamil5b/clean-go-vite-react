package email

import (
	"context"
	"fmt"
	"strconv"
)

// GetEmailLog retrieves an email log entry from GORM
func (r *GORMEmailRepository) GetEmailLog(ctx context.Context, id string) (map[string]interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	numID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("email log not found: %s", id)
	}

	var emailLog EmailLogModel
	if err := r.db.WithContext(ctx).First(&emailLog, uint(numID)).Error; err != nil {
		return nil, fmt.Errorf("email log not found: %s", id)
	}

	return map[string]interface{}{
		"id":      emailLog.ID,
		"to":      emailLog.To,
		"subject": emailLog.Subject,
		"body":    emailLog.Body,
		"status":  emailLog.Status,
	}, nil
}
