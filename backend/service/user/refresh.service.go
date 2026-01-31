package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Refresh generates a new access token from a refresh token
func (s *userService) Refresh(ctx context.Context, refreshToken string) (*response.RefreshResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate refresh token
	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := s.tokenService.GenerateAccessToken(claims.UserID, claims.Email, claims.Name)
	if err != nil {
		return nil, err
	}

	return &response.RefreshResponse{
		Token: accessToken,
	}, nil
}

// GetUser retrieves user information by ID
func (s *userService) GetUser(ctx context.Context, userID string) (*response.GetUser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &response.GetUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
