package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account and returns tokens
func (s *userService) Register(ctx context.Context, req *request.RegisterUserRequest) (*response.RegisterResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Check if user already exists
	existingUser, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user entity
	userEntity := entity.UserEntity{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	// Save to repository
	_, err = s.userRepository.Create(ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateAccessToken(userEntity.ID, userEntity.Email, userEntity.Name)
	if err != nil {
		return nil, err
	}

	// Return response with user and token
	return &response.RegisterResponse{
		User: response.GetUser{
			ID:    userEntity.ID,
			Email: userEntity.Email,
			Name:  userEntity.Name,
		},
		Token: accessToken,
	}, nil
}
