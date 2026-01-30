package message

import (
	"context"
	"testing"
)

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue string
		expectedError bool
	}{
		{
			name:          "should return greeting message successfully",
			expectedValue: "Hello, from the golang World!",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewMessageService()
			result, err := svc.GetMessage(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if result != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, result)
			}
		})
	}
}

func TestGetMessageWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		expectedValue string
		expectedError bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			expectedValue: "Hello, from the golang World!",
			expectedError: false,
		},
		{
			name: "should work with timeout context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx
			},
			expectedValue: "Hello, from the golang World!",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextSetup()
			svc := NewMessageService()
			result, err := svc.GetMessage(ctx)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if result != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, result)
			}
		})
	}
}
