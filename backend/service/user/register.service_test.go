package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name              string
		request           *request.RegisterUserRequest
		mockFindByEmail   *entity.UserEntity
		mockFindByEmailErr error
		mockCreateErr     error
		expectedError     bool
		expectedErrorMsg  string
	}{
		{
			name: "should register user successfully",
			request: &request.RegisterUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockFindByEmail:    nil,
			mockFindByEmailErr: nil,
			mockCreateErr:      nil,
			expectedError:      false,
		},
		{
			name: "should return error when user already exists",
			request: &request.RegisterUserRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Existing User",
			},
			mockFindByEmail: &entity.UserEntity{
				ID:    uuid.New(),
				Email: "existing@example.com",
				Name:  "Existing User",
			},
			mockFindByEmailErr: nil,
			expectedError:      true,
			expectedErrorMsg:   "user already exists",
		},
		{
			name: "should return error when repository find fails",
			request: &request.RegisterUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockFindByEmail:    nil,
			mockFindByEmailErr: errors.New("database error"),
			expectedError:      true,
		},
		{
			name: "should return error when repository create fails",
			request: &request.RegisterUserRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockFindByEmail:    nil,
			mockFindByEmailErr: nil,
			mockCreateErr:      errors.New("create failed"),
			expectedError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := mock.NewMockUserRepository(ctrl)

			// Setup FindByEmail expectation
			mockRepo.EXPECT().
				FindByEmail(gomock.Any(), tt.request.Email).
				Return(tt.mockFindByEmail, tt.mockFindByEmailErr).
				Times(1)

			// Setup Create expectation only if FindByEmail succeeds and user doesn't exist
			if tt.mockFindByEmailErr == nil && tt.mockFindByEmail == nil {
				mockRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&uuid.UUID{}, tt.mockCreateErr).
					Times(1)
			}

			// Setup token service
			tokenConfig := token.TokenConfig{
				AccessTokenSecret:  "test-access-secret",
				RefreshTokenSecret: "test-refresh-secret",
			}
			tokenSvc := token.NewTokenService(tokenConfig)

			// Create service with mocked repository
			svc := NewUserService(mockRepo, tokenSvc)

			// Call the method being tested
			result, err := svc.Register(context.Background(), tt.request)

			// Assert results
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("expected error message '%s', got '%s'", tt.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected non-nil result")
				}
				if result.User.Email != tt.request.Email {
					t.Errorf("expected email %s, got %s", tt.request.Email, result.User.Email)
				}
				if result.User.Name != tt.request.Name {
					t.Errorf("expected name %s, got %s", tt.request.Name, result.User.Name)
				}
				if result.Token == "" {
					t.Errorf("expected non-empty token")
				}
			}
		})
	}
}

func TestRegisterContextCancellation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	tokenConfig := token.TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		RefreshTokenSecret: "test-refresh-secret",
	}
	tokenSvc := token.NewTokenService(tokenConfig)

	svc := NewUserService(mockRepo, tokenSvc)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := &request.RegisterUserRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	result, err := svc.Register(ctx, req)

	if err == nil {
		t.Errorf("expected context.Canceled error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
