package message

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/mock"
)

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name          string
		mockReturn    *string
		mockError     error
		expectedValue string
		expectedError bool
	}{
		{
			name: "should return greeting message successfully",
			mockReturn: func() *string {
				msg := "Hello, from the golang World!"
				return &msg
			}(),
			mockError:     nil,
			expectedValue: "Hello, from the golang World!",
			expectedError: false,
		},
		{
			name: "should return custom message from repository",
			mockReturn: func() *string {
				msg := "Custom greeting"
				return &msg
			}(),
			mockError:     nil,
			expectedValue: "Custom greeting",
			expectedError: false,
		},
		{
			name:          "should return error when repository fails",
			mockReturn:    nil,
			mockError:     errors.New("repository error"),
			expectedValue: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockMessageRepository(ctrl)
			mockRepo.EXPECT().
				GetMessage(gomock.Any(), "default").
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewMessageService(mockRepo)
			result, err := svc.GetMessage(context.Background())

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
				} else if result.Content != tt.expectedValue {
					t.Errorf("expected content %q, got %q", tt.expectedValue, result.Content)
				}
			}
		})
	}
}

func TestGetMessageWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		mockReturn    *string
		mockError     error
		expectedValue string
		expectedError bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			mockReturn: func() *string {
				msg := "Hello, from the golang World!"
				return &msg
			}(),
			mockError:     nil,
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
			mockReturn: func() *string {
				msg := "Hello, from the golang World!"
				return &msg
			}(),
			mockError:     nil,
			expectedValue: "Hello, from the golang World!",
			expectedError: false,
		},
		{
			name: "should handle cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			mockReturn:    nil,
			mockError:     context.Canceled,
			expectedValue: "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockMessageRepository(ctrl)
			ctx := tt.contextSetup()

			mockRepo.EXPECT().
				GetMessage(ctx, "default").
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			svc := NewMessageService(mockRepo)
			result, err := svc.GetMessage(ctx)

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
				} else if result.Content != tt.expectedValue {
					t.Errorf("expected content %q, got %q", tt.expectedValue, result.Content)
				}
			}
		})
	}
}
