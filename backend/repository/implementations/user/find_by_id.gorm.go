package user

import (
	"context"
	"fmt"
	"strconv"
)

// FindByID finds a user by ID in SQLite
func (r *SQLiteUserRepository) FindByID(ctx context.Context, id string) (map[string]interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	numID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	var user UserModel
	if err := r.db.WithContext(ctx).First(&user, uint(numID)).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return map[string]interface{}{
		"id":    fmt.Sprintf("user_%d", user.ID),
		"name":  user.Name,
		"email": user.Email,
	}, nil
}
