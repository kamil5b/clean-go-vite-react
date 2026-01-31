package counter

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// IncrementCounter increments the counter and returns the new value as a response
func (s *counterService) IncrementCounter(ctx context.Context) (*response.GetCounter, error) {
	value, err := s.repo.IncrementCounter(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetCounter{
		Value: value,
	}, nil
}
