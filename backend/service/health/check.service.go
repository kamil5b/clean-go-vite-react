package health

import (
	"context"
)

// Check performs a basic health check
func (s *healthService) Check(ctx context.Context) (*HealthStatus, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return &HealthStatus{
		Status:  "ok",
		Message: "Service is healthy",
		Details: map[string]interface{}{
			"timestamp": ctx.Value("timestamp"),
		},
	}, nil
}
