package email

import (
	"context"
	"testing"
)

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name          string
		to            string
		subject       string
		body          string
		expectedError bool
	}{
		{
			name:          "should send email successfully",
			to:            "test@example.com",
			subject:       "Test Subject",
			body:          "Test body content",
			expectedError: false,
		},
		{
			name:          "should handle multiple recipients",
			to:            "user1@example.com",
			subject:       "Multi recipient",
			body:          "Body",
			expectedError: false,
		},
		{
			name:          "should handle empty body",
			to:            "test@example.com",
			subject:       "Subject",
			body:          "",
			expectedError: false,
		},
		{
			name:          "should handle special characters in subject",
			to:            "test@example.com",
			subject:       "Subject with special chars: !@#$%",
			body:          "Content",
			expectedError: false,
		},
		{
			name:          "should handle long email content",
			to:            "test@example.com",
			subject:       "Long content",
			body:          "This is a very long email body with lots of content that should still be handled correctly without any issues or errors",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService()
			err := svc.SendEmail(context.Background(), tt.to, tt.subject, tt.body)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSendEmailWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		to            string
		subject       string
		body          string
		expectedError bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			to:            "test@example.com",
			subject:       "Subject",
			body:          "Body",
			expectedError: false,
		},
		{
			name: "should handle cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			to:            "test@example.com",
			subject:       "Subject",
			body:          "Body",
			expectedError: true,
		},
		{
			name: "should handle deadline exceeded context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			to:            "test@example.com",
			subject:       "Subject",
			body:          "Body",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService()
			ctx := tt.contextSetup()
			err := svc.SendEmail(ctx, tt.to, tt.subject, tt.body)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
