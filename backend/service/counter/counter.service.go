package counter

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
)

// CounterService defines the interface for counter operations
type CounterService interface {
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)
}

// counterService is the concrete implementation of CounterService
type counterService struct {
	repo interfaces.CounterRepository
}

// NewCounterService creates a new instance of CounterService
func NewCounterService(repo interfaces.CounterRepository) CounterService {
	return &counterService{
		repo: repo,
	}
}
