package user

import (
	"context"
	"errors"

	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates a user and returns tokens
func (s *userService) Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Find user by email
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateAccessToken(user.ID, user.Email, user.Name)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Store refresh token in database/cache for revocation
	_ = refreshToken

	return &response.LoginResponse{
		User: response.GetUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Token: accessToken,
	}, nil
}
