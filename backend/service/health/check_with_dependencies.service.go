package health

import (
	"context"
	"fmt"
)

// CheckWithDependencies performs health checks on multiple dependencies
func (s *healthService) CheckWithDependencies(ctx context.Context, checks map[string]func(context.Context) error) (*HealthStatus, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	status := &HealthStatus{
		Status:  "ok",
		Message: "All dependencies healthy",
		Details: make(map[string]interface{}),
	}

	unhealthyCount := 0
	for name, check := range checks {
		if err := check(ctx); err != nil {
			status.Details[name] = map[string]string{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			unhealthyCount++
		} else {
			status.Details[name] = map[string]string{
				"status": "healthy",
			}
		}
	}

	if unhealthyCount > 0 {
		status.Status = "degraded"
		status.Message = fmt.Sprintf("%d dependency checks failed", unhealthyCount)
	}

	return status, nil
}
