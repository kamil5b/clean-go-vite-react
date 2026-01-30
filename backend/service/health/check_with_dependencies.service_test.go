package health

import (
	"context"
	"errors"
	"testing"
)

func TestCheckWithDependencies(t *testing.T) {
	tests := []struct {
		name            string
		checks          map[string]func(context.Context) error
		expectedStatus  string
		expectedMessage string
		expectedError   bool
		healthyCount    int
		unhealthyCount  int
	}{
		{
			name: "should return ok when all dependencies are healthy",
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return nil
				},
				"cache": func(ctx context.Context) error {
					return nil
				},
			},
			expectedStatus:  "ok",
			expectedMessage: "All dependencies healthy",
			expectedError:   false,
			healthyCount:    2,
			unhealthyCount:  0,
		},
		{
			name: "should return degraded when some dependencies are unhealthy",
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return nil
				},
				"cache": func(ctx context.Context) error {
					return errors.New("connection timeout")
				},
			},
			expectedStatus:  "degraded",
			expectedMessage: "1 dependency checks failed",
			expectedError:   false,
			healthyCount:    1,
			unhealthyCount:  1,
		},
		{
			name: "should return degraded when all dependencies are unhealthy",
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return errors.New("connection refused")
				},
				"cache": func(ctx context.Context) error {
					return errors.New("timeout")
				},
				"queue": func(ctx context.Context) error {
					return errors.New("unavailable")
				},
			},
			expectedStatus:  "degraded",
			expectedMessage: "3 dependency checks failed",
			expectedError:   false,
			healthyCount:    0,
			unhealthyCount:  3,
		},
		{
			name:            "should handle empty checks",
			checks:          map[string]func(context.Context) error{},
			expectedStatus:  "ok",
			expectedMessage: "All dependencies healthy",
			expectedError:   false,
			healthyCount:    0,
			unhealthyCount:  0,
		},
		{
			name: "should handle single healthy dependency",
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return nil
				},
			},
			expectedStatus:  "ok",
			expectedMessage: "All dependencies healthy",
			expectedError:   false,
			healthyCount:    1,
			unhealthyCount:  0,
		},
		{
			name: "should handle single unhealthy dependency",
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return errors.New("connection failed")
				},
			},
			expectedStatus:  "degraded",
			expectedMessage: "1 dependency checks failed",
			expectedError:   false,
			healthyCount:    0,
			unhealthyCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewHealthService()
			status, err := svc.CheckWithDependencies(context.Background(), tt.checks)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if status == nil {
				t.Fatalf("expected non-nil status")
			}

			if status.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, status.Status)
			}

			if status.Message != tt.expectedMessage {
				t.Errorf("expected message %s, got %s", tt.expectedMessage, status.Message)
			}

			if status.Details == nil {
				t.Errorf("expected non-nil details")
			}
		})
	}
}

func TestCheckWithDependenciesWithContext(t *testing.T) {
	tests := []struct {
		name           string
		contextSetup   func() context.Context
		checks         map[string]func(context.Context) error
		expectedStatus string
		expectedError  bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return nil
				},
			},
			expectedStatus: "ok",
			expectedError:  false,
		},
		{
			name: "should work with cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			checks: map[string]func(context.Context) error{
				"database": func(ctx context.Context) error {
					return ctx.Err()
				},
			},
			expectedStatus: "degraded",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.contextSetup()
			svc := NewHealthService()
			status, err := svc.CheckWithDependencies(ctx, tt.checks)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if status == nil && !tt.expectedError {
				t.Errorf("expected non-nil status")
			}
		})
	}
}
