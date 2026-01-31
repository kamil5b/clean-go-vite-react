package counter

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// GetCounter returns the current counter value using the response model
func (s *counterService) GetCounter(ctx context.Context) (*response.GetCounter, error) {
	value, err := s.repo.GetCounter(ctx)
	if err != nil {
		return nil, err
	}
	return &response.GetCounter{
		Value: value,
	}, nil
}
