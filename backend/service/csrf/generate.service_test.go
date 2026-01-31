package csrf

import (
	"encoding/hex"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "should generate csrf token successfully",
			expectedError: false,
		},
		{
			name:          "should generate different tokens on each call",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewCSRFService()
			token, err := service.GenerateToken()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token == "" {
					t.Errorf("expected non-empty token, got empty string")
				}

				// Verify token is valid hex
				_, hexErr := hex.DecodeString(token)
				if hexErr != nil {
					t.Errorf("token is not valid hex: %v", hexErr)
				}

				// Verify token length (32 bytes = 64 hex characters)
				if len(token) != 64 {
					t.Errorf("expected token length 64, got %d", len(token))
				}
			}
		})
	}
}

func TestGenerateTokenUniqueness(t *testing.T) {
	service := NewCSRFService()

	token1, err1 := service.GenerateToken()
	if err1 != nil {
		t.Fatalf("first token generation failed: %v", err1)
	}

	token2, err2 := service.GenerateToken()
	if err2 != nil {
		t.Fatalf("second token generation failed: %v", err2)
	}

	if token1 == token2 {
		t.Errorf("expected different tokens, but got same token: %s", token1)
	}
}

func TestGenerateTokenFormat(t *testing.T) {
	service := NewCSRFService()

	for i := 0; i < 10; i++ {
		token, err := service.GenerateToken()

		if err != nil {
			t.Fatalf("token generation failed: %v", err)
		}

		// Verify token is valid hex
		_, hexErr := hex.DecodeString(token)
		if hexErr != nil {
			t.Errorf("token is not valid hex: %v", hexErr)
		}

		// Verify token is exactly 64 characters (32 bytes in hex)
		if len(token) != 64 {
			t.Errorf("expected token length 64, got %d", len(token))
		}

		// Verify token contains only hex characters
		for _, ch := range token {
			if !isHexChar(rune(ch)) {
				t.Errorf("token contains non-hex character: %c", ch)
			}
		}
	}
}

func isHexChar(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func TestGenerateTokenRandomness(t *testing.T) {
	service := NewCSRFService()
	tokens := make(map[string]bool)

	// Generate 100 tokens and verify they're all unique
	for i := 0; i < 100; i++ {
		token, err := service.GenerateToken()

		if err != nil {
			t.Fatalf("token generation failed on iteration %d: %v", i, err)
		}

		if tokens[token] {
			t.Errorf("generated duplicate token on iteration %d: %s", i, token)
		}

		tokens[token] = true
	}

	if len(tokens) != 100 {
		t.Errorf("expected 100 unique tokens, got %d", len(tokens))
	}
}

func TestGenerateTokenNonEmpty(t *testing.T) {
	service := NewCSRFService()
	token, err := service.GenerateToken()

	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Errorf("generated token is empty")
	}

	if len(token) == 0 {
		t.Errorf("generated token has zero length")
	}
}
