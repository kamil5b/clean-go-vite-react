package user

import (
	"context"
	"fmt"
	"strconv"
)

// Update updates a user in SQLite
func (r *SQLiteUserRepository) Update(ctx context.Context, id string, user map[string]interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	numID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("user not found: %s", id)
	}

	updateData := map[string]interface{}{
		"name":  user["name"],
		"email": user["email"],
	}

	if err := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", uint(numID)).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
