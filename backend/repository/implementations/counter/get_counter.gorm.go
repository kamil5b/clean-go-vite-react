package counter

import "context"

// GetCounter returns the current counter value from SQLite
func (r *SQLiteCounterRepository) GetCounter(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	var counter CounterModel
	if err := r.db.WithContext(ctx).First(&counter).Error; err != nil {
		return 0, err
	}

	return counter.Value, nil
}
