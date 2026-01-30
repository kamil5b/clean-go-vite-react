package health

import (
	"context"
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
