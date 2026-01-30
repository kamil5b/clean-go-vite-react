package health

import (
	"context"
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name            string
		contextSetup    func() context.Context
		expectedStatus  string
		expectedMessage string
		expectedError   bool
	}{
		{
			name: "should return ok status on successful check",
			contextSetup: func() context.Context {
				return context.Background()
			},
			expectedStatus:  "ok",
			expectedMessage: "Service is healthy",
			expectedError:   false,
		},
		{
			name: "should handle context with value",
			contextSetup: func() context.Context {
				return context.WithValue(context.Background(), "timestamp", "2024-01-01T00:00:00Z")
			},
			expectedStatus:  "ok",
			expectedMessage: "Service is healthy",
			expectedError:   false,
		},
		{
			name: "should return error on cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expectedStatus:  "",
			expectedMessage: "",
			expectedError:   true,
		},
		{
			name: "should return error on deadline exceeded",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 0)
				defer cancel()
				return ctx
			},
			expectedStatus:  "",
			expectedMessage: "",
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewHealthService()
			ctx := tt.contextSetup()

			result, err := svc.Check(ctx)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !tt.expectedError && result != nil {
				if result.Status != tt.expectedStatus {
					t.Errorf("expected status %s, got %s", tt.expectedStatus, result.Status)
				}
				if result.Message != tt.expectedMessage {
					t.Errorf("expected message %s, got %s", tt.expectedMessage, result.Message)
				}
			}
		})
	}
}

func TestCheckMultipleCalls(t *testing.T) {
	tests := []struct {
		name      string
		callCount int
	}{
		{
			name:      "should handle single call",
			callCount: 1,
		},
		{
			name:      "should handle multiple sequential calls",
			callCount: 5,
		},
		{
			name:      "should handle concurrent-like calls",
			callCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewHealthService()

			for i := 0; i < tt.callCount; i++ {
				result, err := svc.Check(context.Background())
				if err != nil {
					t.Errorf("call %d: unexpected error: %v", i+1, err)
				}
				if result == nil {
					t.Errorf("call %d: expected non-nil result", i+1)
				}
				if result.Status != "ok" {
					t.Errorf("call %d: expected status ok, got %s", i+1, result.Status)
				}
			}
		})
	}
}
