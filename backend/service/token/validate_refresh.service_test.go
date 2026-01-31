package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateRefreshToken(t *testing.T) {
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
	}{
		{
			name: "should validate refresh token successfully",
			tokenFunc: func(svc TokenService) string {
				token, _ := svc.GenerateRefreshToken(testUserID)
				return token
			},
			expectedError: false,
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
				token, _ := svc.GenerateRefreshToken(testUserID)
				return token + "tampered"
			},
			expectedError: true,
		},
		{
			name: "should return error for wrong secret",
			tokenFunc: func(svc TokenService) string {
				wrongConfig := TokenConfig{
					AccessTokenSecret:  "test-access-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenSecret: "wrong-refresh-secret",
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				}
				wrongSvc := NewTokenService(wrongConfig)
				token, _ := wrongSvc.GenerateRefreshToken(testUserID)
				return token
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTokenService(config)
			tokenString := tt.tokenFunc(service)

			claims, err := service.ValidateRefreshToken(tokenString)

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
				}
			}
		})
	}
}

func TestValidateRefreshTokenExpired(t *testing.T) {
	testUserID := uuid.New()
	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: -1 * time.Second, // Already expired
	}

	service := NewTokenService(config)
	tokenString, _ := service.GenerateRefreshToken(testUserID)

	// Wait a moment to ensure token is definitely expired
	time.Sleep(100 * time.Millisecond)

	claims, err := service.ValidateRefreshToken(tokenString)

	if err == nil {
		t.Errorf("expected error for expired token, got nil")
	}
	if claims != nil {
		t.Errorf("expected nil claims for invalid token")
	}
}



func TestValidateRefreshTokenClaimsIntegrity(t *testing.T) {
	testUserID := uuid.New()

	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	service := NewTokenService(config)
	tokenString, _ := service.GenerateRefreshToken(testUserID)

	claims, err := service.ValidateRefreshToken(tokenString)

	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != testUserID {
		t.Errorf("UserID mismatch: expected %v, got %v", testUserID, claims.UserID)
	}
	if claims.Issuer != "go-vite-react" {
		t.Errorf("Issuer mismatch: expected go-vite-react, got %s", claims.Issuer)
	}
	// Refresh token should not have email or name
	if claims.Email != "" {
		t.Errorf("refresh token should not have email, got %s", claims.Email)
	}
	if claims.Name != "" {
		t.Errorf("refresh token should not have name, got %s", claims.Name)
	}
}

func TestValidateRefreshTokenAccessTokenMismatch(t *testing.T) {
	testUserID := uuid.New()
	config := TokenConfig{
		AccessTokenSecret:  "test-access-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenSecret: "test-refresh-secret",
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	service := NewTokenService(config)

	// Generate an access token
	accessToken, _ := service.GenerateAccessToken(testUserID, "test@example.com", "Test User")

	// Try to validate it as a refresh token (with wrong secret)
	_, err := service.ValidateRefreshToken(accessToken)

	if err == nil {
		t.Errorf("expected error when validating access token as refresh token, got nil")
	}
}
