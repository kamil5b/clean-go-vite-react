package counter

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
)

func TestGetCounter(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    int
		mockError     error
		expectedValue int
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "should return counter value successfully",
			mockReturn:    42,
			mockError:     nil,
			expectedValue: 42,
			expectedError: false,
		},
		{
			name:          "should return zero value successfully",
			mockReturn:    0,
			mockError:     nil,
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "should handle negative counter value",
			mockReturn:    -10,
			mockError:     nil,
			expectedValue: -10,
			expectedError: false,
		},
		{
			name:          "should handle large counter value",
			mockReturn:    999999999,
			mockError:     nil,
			expectedValue: 999999999,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    0,
			mockError:     errors.New("database connection failed"),
			expectedValue: 0,
			expectedError: true,
			errorMessage:  "database connection failed",
		},
		{
			name:          "should return error on context canceled",
			mockReturn:    0,
			mockError:     context.Canceled,
			expectedValue: 0,
			expectedError: true,
			errorMessage:  "context canceled",
		},
		{
			name:          "should return error on context deadline exceeded",
			mockReturn:    0,
			mockError:     context.DeadlineExceeded,
			expectedValue: 0,
			expectedError: true,
			errorMessage:  "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockCounterRepository(ctrl)
			mockRepo.EXPECT().
				GetCounter(gomock.Any()).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewCounterService(mockRepo)
			result, err := svc.GetCounter(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !tt.expectedError {
				if result == nil {
					t.Errorf("expected result, got nil")
				} else if result.Value != tt.expectedValue {
					t.Errorf("expected value %d, got %d", tt.expectedValue, result.Value)
				}
			}
		})
	}
}

func TestGetCounterWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		mockReturn    int
		mockError     error
		expectedValue int
		expectedError bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			mockReturn:    10,
			mockError:     nil,
			expectedValue: 10,
			expectedError: false,
		},
		{
			name: "should work with cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			mockReturn:    0,
			mockError:     context.Canceled,
			expectedValue: 0,
			expectedError: true,
		},
		{
			name: "should pass context to repository",
			contextSetup: func() context.Context {
				return context.Background()
			},
			mockReturn:    25,
			mockError:     nil,
			expectedValue: 25,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockCounterRepository(ctrl)
			ctx := tt.contextSetup()

			mockRepo.EXPECT().
				GetCounter(ctx).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewCounterService(mockRepo)
			result, err := svc.GetCounter(ctx)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !tt.expectedError {
				if result == nil {
					t.Errorf("expected result, got nil")
				} else if result.Value != tt.expectedValue {
					t.Errorf("expected value %d, got %d", tt.expectedValue, result.Value)
				}
			}
		})
	}
}
