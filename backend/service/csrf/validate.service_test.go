package csrf

import (
	"testing"
)

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		expectedValid  bool
	}{
		{
			name:          "should validate non-empty token",
			token:         "valid-token-123",
			expectedValid: true,
		},
		{
			name:          "should validate token with special characters",
			token:         "token!@#$%^&*()",
			expectedValid: true,
		},
		{
			name:          "should validate single character token",
			token:         "a",
			expectedValid: true,
		},
		{
			name:          "should reject empty token",
			token:         "",
			expectedValid: false,
		},
		{
			name:          "should validate hex token",
			token:         "a1b2c3d4e5f6",
			expectedValid: true,
		},
		{
			name:          "should validate long token",
			token:         "thisisaverylongtokenthatcontainsmanycharactersandshouldbeconsideredvalid",
			expectedValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewCSRFService()
			result := service.ValidateToken(tt.token)

			if result != tt.expectedValid {
				t.Errorf("expected validation result %v, got %v for token %q", tt.expectedValid, result, tt.token)
			}
		})
	}
}

func TestValidateTokenEmptyString(t *testing.T) {
	service := NewCSRFService()
	result := service.ValidateToken("")

	if result {
		t.Errorf("expected empty token to be invalid, got valid")
	}
}

func TestValidateTokenNonEmpty(t *testing.T) {
	service := NewCSRFService()
	testTokens := []string{
		"token",
		"a",
		"1",
		"!@#$%",
		"very-long-token-with-many-characters-12345",
	}

	for _, token := range testTokens {
		result := service.ValidateToken(token)
		if !result {
			t.Errorf("expected token %q to be valid, got invalid", token)
		}
	}
}

func TestValidateTokenWithGeneratedToken(t *testing.T) {
	service := NewCSRFService()

	generatedToken, err := service.GenerateToken()
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	result := service.ValidateToken(generatedToken)
	if !result {
		t.Errorf("expected generated token to be valid, got invalid")
	}
}

func TestValidateTokenConsistency(t *testing.T) {
	service := NewCSRFService()
	token := "consistent-token"

	result1 := service.ValidateToken(token)
	result2 := service.ValidateToken(token)
	result3 := service.ValidateToken(token)

	if result1 != result2 || result2 != result3 {
		t.Errorf("validation should be consistent for same token")
	}
}

func TestValidateTokenWhitespace(t *testing.T) {
	service := NewCSRFService()

	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "should validate token with spaces",
			token:    "token with spaces",
			expected: true,
		},
		{
			name:     "should validate token with tabs",
			token:    "token\twith\ttabs",
			expected: true,
		},
		{
			name:     "should validate token with newlines",
			token:    "token\nwith\nnewlines",
			expected: true,
		},
		{
			name:     "should not validate only spaces",
			token:    "   ",
			expected: true, // Spaces are still non-empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ValidateToken(tt.token)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for token %q", tt.expected, result, tt.token)
			}
		})
	}
}
