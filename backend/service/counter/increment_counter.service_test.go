package counter

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
)

func TestIncrementCounter(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    int
		mockError     error
		expectedValue int
		expectedError bool
		errorMessage  string
	}{
		{
			name:          "should increment counter successfully",
			mockReturn:    43,
			mockError:     nil,
			expectedValue: 43,
			expectedError: false,
		},
		{
			name:          "should increment from zero",
			mockReturn:    1,
			mockError:     nil,
			expectedValue: 1,
			expectedError: false,
		},
		{
			name:          "should handle large incremented value",
			mockReturn:    1000000,
			mockError:     nil,
			expectedValue: 1000000,
			expectedError: false,
		},
		{
			name:          "should handle negative counter after increment",
			mockReturn:    -5,
			mockError:     nil,
			expectedValue: -5,
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    0,
			mockError:     errors.New("failed to increment counter"),
			expectedValue: 0,
			expectedError: true,
			errorMessage:  "failed to increment counter",
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
		{
			name:          "should return error on database error",
			mockReturn:    0,
			mockError:     errors.New("database unavailable"),
			expectedValue: 0,
			expectedError: true,
			errorMessage:  "database unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockCounterRepository(ctrl)
			mockRepo.EXPECT().
				IncrementCounter(gomock.Any()).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewCounterService(mockRepo)
			result, err := svc.IncrementCounter(context.Background())

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected non-nil result, got nil")
				} else if result.Value != tt.expectedValue {
					t.Errorf("expected value %d, got %d", tt.expectedValue, result.Value)
				}
			}
		})
	}
}

func TestIncrementCounterWithContext(t *testing.T) {
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
			mockReturn:    50,
			mockError:     nil,
			expectedValue: 50,
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
			mockReturn:    77,
			mockError:     nil,
			expectedValue: 77,
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
				IncrementCounter(ctx).
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewCounterService(mockRepo)
			result, err := svc.IncrementCounter(ctx)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected non-nil result, got nil")
				} else if result.Value != tt.expectedValue {
					t.Errorf("expected value %d, got %d", tt.expectedValue, result.Value)
				}
			}
		})
	}
}

func TestIncrementCounterSequential(t *testing.T) {
	tests := []struct {
		name          string
		callCount     int
		startValue    int
		expectedCalls []int
	}{
		{
			name:          "should handle multiple sequential increments",
			callCount:     3,
			startValue:    1,
			expectedCalls: []int{1, 2, 3},
		},
		{
			name:          "should handle single increment",
			callCount:     1,
			startValue:    5,
			expectedCalls: []int{5},
		},
		{
			name:          "should handle multiple increments from zero",
			callCount:     5,
			startValue:    0,
			expectedCalls: []int{0, 1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockCounterRepository(ctrl)

			for i := 0; i < tt.callCount; i++ {
				mockRepo.EXPECT().
					IncrementCounter(gomock.Any()).
					Return(tt.expectedCalls[i], nil)
			}

			svc := NewCounterService(mockRepo)

			for i := 0; i < tt.callCount; i++ {
				result, err := svc.IncrementCounter(context.Background())
				if err != nil {
					t.Errorf("call %d: unexpected error: %v", i+1, err)
				}
				if result == nil {
					t.Errorf("call %d: expected non-nil result, got nil", i+1)
				} else if result.Value != tt.expectedCalls[i] {
					t.Errorf("call %d: expected %d, got %d", i+1, tt.expectedCalls[i], result.Value)
				}
			}
		})
	}
}
