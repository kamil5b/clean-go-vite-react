package counter

import "context"

// GetCounter returns the current counter value
func (s *counterService) GetCounter(ctx context.Context) (int, error) {
	return s.repo.GetCounter(ctx)
}
