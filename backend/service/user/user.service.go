package user

import (
	"context"

	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
)

// UserService defines the interface for user operations
type UserService interface {
	Register(ctx context.Context, req *request.RegisterUserRequest) (*response.RegisterResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.LoginResponse, error)
	Refresh(ctx context.Context, userID string) (*response.RefreshResponse, error)
	GetUser(ctx context.Context, userID string) (*response.GetUser, error)
}

// userService is the concrete implementation of UserService
type userService struct {
	userRepository interfaces.UserRepository
	tokenService   token.TokenService
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepository interfaces.UserRepository, tokenService token.TokenService) UserService {
	return &userService{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}
