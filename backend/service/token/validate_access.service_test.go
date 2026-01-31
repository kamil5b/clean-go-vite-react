package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateAccessToken(t *testing.T) {
	testUserID := uuid.New()
	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	tests := []struct {
		name          string
		tokenFunc     func(TokenService) string
		expectedError bool
		errorMsg      string
		validateEmail string
	}{
		{
			name: "should validate access token successfully",
			tokenFunc: func(svc TokenService) string {
				token, _ := svc.GenerateAccessToken(testUserID, "test@example.com", "Test User")
				return token
			},
			expectedError: false,
			validateEmail: "test@example.com",
		},
		{
			name: "should return error for empty token",
			tokenFunc: func(svc TokenService) string {
				return ""
			},
			expectedError: true,
			errorMsg:      "token is empty",
		},
		{
			name: "should return error for invalid token format",
			tokenFunc: func(svc TokenService) string {
				return "invalid.token.format"
			},
			expectedError: true,
		},
		{
			name: "should return error for tampered token",
			tokenFunc: func(svc TokenService) string {
				token, _ := svc.GenerateAccessToken(testUserID, "test@example.com", "Test User")
				return token + "tampered"
			},
			expectedError: true,
		},
		{
			name: "should return error for wrong secret",
			tokenFunc: func(svc TokenService) string {
				wrongConfig := TokenConfig{
					AccessTokenSecret:  "wrong-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenSecret: "test-refresh-secret",
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				}
				wrongSvc := NewTokenService(wrongConfig)
				token, _ := wrongSvc.GenerateAccessToken(testUserID, "test@example.com", "Test User")
				return token
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTokenService(config)
			tokenString := tt.tokenFunc(service)

			claims, err := service.ValidateAccessToken(tokenString)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if claims == nil {
					t.Errorf("expected non-nil claims")
				}
				if claims != nil {
					if claims.UserID != testUserID {
						t.Errorf("expected UserID %v, got %v", testUserID, claims.UserID)
					}
					if claims.Email != tt.validateEmail {
						t.Errorf("expected Email %s, got %s", tt.validateEmail, claims.Email)
					}
				}
			}
		})
	}
}

func TestValidateAccessTokenExpired(t *testing.T) {
	testUserID := uuid.New()
	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  -1 * time.Second,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	service := NewTokenService(config)
	tokenString, _ := service.GenerateAccessToken(testUserID, "test@example.com", "Test User")

	time.Sleep(100 * time.Millisecond)

	claims, err := service.ValidateAccessToken(tokenString)

	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
	if claims != nil {
		t.Errorf("expected nil claims for invalid token")
	}
}

func TestValidateAccessTokenClaimsIntegrity(t *testing.T) {
	testUserID := uuid.New()
	testEmail := "test@example.com"
	testName := "Test User"

	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	service := NewTokenService(config)
	tokenString, _ := service.GenerateAccessToken(testUserID, testEmail, testName)

	claims, err := service.ValidateAccessToken(tokenString)

	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != testUserID {
		t.Errorf("UserID mismatch: expected %v, got %v", testUserID, claims.UserID)
	}
	if claims.Email != testEmail {
		t.Errorf("Email mismatch: expected %s, got %s", testEmail, claims.Email)
	}
	if claims.Name != testName {
		t.Errorf("Name mismatch: expected %s, got %s", testName, claims.Name)
	}
	if claims.Issuer != "go-vite-react" {
		t.Errorf("Issuer mismatch: expected go-vite-react, got %s", claims.Issuer)
	}
}
