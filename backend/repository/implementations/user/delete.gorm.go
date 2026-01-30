package user

import (
	"context"
	"strconv"
)

// Delete deletes a user in SQLite
func (r *SQLiteUserRepository) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	numID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&UserModel{}, uint(numID)).Error; err != nil {
		return err
	}

	return nil
}
