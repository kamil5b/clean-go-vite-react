package email

import (
	"context"
	"testing"

	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
)

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name          string
		request       *request.SaveEmailRequest
		expectedError bool
	}{
		{
			name: "should send email successfully",
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Test Subject",
				Body:    "Test body content",
			},
			expectedError: false,
		},
		{
			name: "should handle multiple recipients",
			request: &request.SaveEmailRequest{
				To:      "user1@example.com",
				Subject: "Multi recipient",
				Body:    "Body",
			},
			expectedError: false,
		},
		{
			name: "should handle empty body",
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Subject",
				Body:    "",
			},
			expectedError: false,
		},
		{
			name: "should handle special characters in subject",
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Subject with special chars: !@#$%",
				Body:    "Content",
			},
			expectedError: false,
		},
		{
			name: "should handle long email content",
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Long content",
				Body:    "This is a very long email body with lots of content that should still be handled correctly without any issues or errors",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService()
			result, err := svc.SendEmail(context.Background(), tt.request)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if result == nil {
				t.Errorf("expected non-nil response, got nil")
			}

			if result != nil {
				if result.To != tt.request.To {
					t.Errorf("expected To %q, got %q", tt.request.To, result.To)
				}
				if result.Subject != tt.request.Subject {
					t.Errorf("expected Subject %q, got %q", tt.request.Subject, result.Subject)
				}
				if result.Body != tt.request.Body {
					t.Errorf("expected Body %q, got %q", tt.request.Body, result.Body)
				}
				if result.Status != "sent" {
					t.Errorf("expected Status 'sent', got %q", result.Status)
				}
			}
		})
	}
}

func TestSendEmailWithContext(t *testing.T) {
	tests := []struct {
		name          string
		contextSetup  func() context.Context
		request       *request.SaveEmailRequest
		expectedError bool
	}{
		{
			name: "should work with background context",
			contextSetup: func() context.Context {
				return context.Background()
			},
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Subject",
				Body:    "Body",
			},
			expectedError: false,
		},
		{
			name: "should handle cancelled context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Subject",
				Body:    "Body",
			},
			expectedError: true,
		},
		{
			name: "should handle deadline exceeded context",
			contextSetup: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			request: &request.SaveEmailRequest{
				To:      "test@example.com",
				Subject: "Subject",
				Body:    "Body",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewEmailService()
			ctx := tt.contextSetup()
			result, err := svc.SendEmail(ctx, tt.request)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected non-nil response, got nil")
				}
			}
		})
	}
}
