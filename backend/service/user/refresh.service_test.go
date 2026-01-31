package user

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
)

func TestRefresh(t *testing.T) {
	testUserID := uuid.New()
	tokenConfig := token.TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		RefreshTokenSecret: "test-refresh-secret",
	}
	tokenSvc := token.NewTokenService(tokenConfig)

	// Generate a valid refresh token for testing
	validRefreshToken, _ := tokenSvc.GenerateRefreshToken(testUserID)

	tests := []struct {
		name             string
		refreshToken     string
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:          "should refresh token successfully",
			refreshToken:  validRefreshToken,
			expectedError: false,
		},
		{
			name:             "should return error when refresh token is invalid",
			refreshToken:     "invalid-token",
			expectedError:    true,
			expectedErrorMsg: "invalid refresh token",
		},
		{
			name:             "should return error when refresh token is empty",
			refreshToken:     "",
			expectedError:    true,
			expectedErrorMsg: "invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := mock.NewMockUserRepository(ctrl)

			// Create service with same token config
			svc := NewUserService(mockRepo, tokenSvc)

			// Call refresh
			result, err := svc.Refresh(context.Background(), tt.refreshToken)

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
				if result != nil && result.Token == "" {
					t.Errorf("expected non-empty token")
				}
			}
		})
	}
}

func TestRefreshContextCancellation(t *testing.T) {
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

	result, err := svc.Refresh(ctx, "some-token")

	if err == nil {
		t.Errorf("expected context.Canceled error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetUser(t *testing.T) {
	testUserID := uuid.New()

	tests := []struct {
		name             string
		userIDString     string
		mockFindByID     *entity.UserEntity
		mockFindByIDErr  error
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:         "should get user successfully",
			userIDString: testUserID.String(),
			mockFindByID: &entity.UserEntity{
				ID:    testUserID,
				Email: "test@example.com",
				Name:  "Test User",
			},
			mockFindByIDErr: nil,
			expectedError:   false,
		},
		{
			name:             "should return error when user id is invalid",
			userIDString:     "invalid-uuid",
			mockFindByID:     nil,
			expectedError:    true,
			expectedErrorMsg: "invalid user id",
		},
		{
			name:             "should return error when user not found",
			userIDString:     testUserID.String(),
			mockFindByID:     nil,
			mockFindByIDErr:  nil,
			expectedError:    true,
			expectedErrorMsg: "user not found",
		},
		{
			name:            "should return error when repository fails",
			userIDString:    testUserID.String(),
			mockFindByID:    nil,
			mockFindByIDErr: errors.New("database error"),
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Setup mock repository
			mockRepo := mock.NewMockUserRepository(ctrl)

			// Setup FindByID expectation only if valid UUID
			if _, err := uuid.Parse(tt.userIDString); err == nil {
				mockRepo.EXPECT().
					FindByID(gomock.Any(), gomock.Any()).
					Return(tt.mockFindByID, tt.mockFindByIDErr).
					Times(1)
			}

			// Setup token service
			tokenConfig := token.TokenConfig{
				AccessTokenSecret:  "test-access-secret",
				RefreshTokenSecret: "test-refresh-secret",
			}
			tokenSvc := token.NewTokenService(tokenConfig)

			// Create service
			svc := NewUserService(mockRepo, tokenSvc)

			// Call getuser
			result, err := svc.GetUser(context.Background(), tt.userIDString)

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
			}
		})
	}
}

func TestGetUserContextCancellation(t *testing.T) {
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

	result, err := svc.GetUser(ctx, uuid.New().String())

	if err == nil {
		t.Errorf("expected context.Canceled error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
