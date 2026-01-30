package counter

import "context"

// IncrementCounter increments the counter and returns the new value
func (s *counterService) IncrementCounter(ctx context.Context) (int, error) {
	return s.repo.IncrementCounter(ctx)
}
