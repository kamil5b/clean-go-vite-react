package service

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/internal/repository"
)

// CounterService defines the interface for counter operations
type CounterService interface {
	GetCounter(ctx context.Context) (int, error)
	IncrementCounter(ctx context.Context) (int, error)
}

// counterService is the concrete implementation of CounterService
type counterService struct {
	repo repository.CounterRepository
}

// NewCounterService creates a new instance of CounterService
func NewCounterService(repo repository.CounterRepository) CounterService {
	return &counterService{
		repo: repo,
	}
}

// GetCounter returns the current counter value
func (s *counterService) GetCounter(ctx context.Context) (int, error) {
	return s.repo.GetCounter(ctx)
}

// IncrementCounter increments the counter and returns the new value
func (s *counterService) IncrementCounter(ctx context.Context) (int, error) {
	return s.repo.IncrementCounter(ctx)
}
