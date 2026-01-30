package service

import (
	"context"
	"fmt"
)

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// HealthService defines the interface for health checks
type HealthService interface {
	Check(ctx context.Context) (*HealthStatus, error)
	CheckWithDependencies(ctx context.Context, checks map[string]func(context.Context) error) (*HealthStatus, error)
}

// healthService is the concrete implementation of HealthService
type healthService struct{}

// NewHealthService creates a new instance of HealthService
func NewHealthService() HealthService {
	return &healthService{}
}

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
