package user

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	tokenConfig := token.TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		RefreshTokenSecret: "test-refresh-secret",
	}
	tokenSvc := token.NewTokenService(tokenConfig)

	service := NewUserService(mockRepo, tokenSvc)

	if service == nil {
		t.Errorf("expected non-nil service, got nil")
	}

	// Verify service is a UserService
	_, ok := service.(*userService)
	if !ok {
		t.Errorf("expected *userService, got %T", service)
	}
}
