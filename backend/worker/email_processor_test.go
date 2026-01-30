package worker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/kamil5b/clean-go-vite-react/backend/task"
)

// MockEmailService is a mock implementation of EmailService for testing
type MockEmailService struct {
	SendEmailFunc func(ctx context.Context, to, subject, body string) error
}

func (m *MockEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(ctx, to, subject, body)
	}
	return nil
}

func TestEmailProcessor_ProcessTask_Success(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			return nil
		},
	}
	processor := NewEmailProcessor(mockService)

	payload := &task.EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	payloadBytes, _ := payload.Marshal()

	task := asynq.NewTask(task.TypeEmailNotification, payloadBytes)

	// Act
	err := processor.ProcessTask(context.Background(), task)

	// Assert
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}
}

func TestEmailProcessor_ProcessTask_InvalidPayload(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			return nil
		},
	}
	processor := NewEmailProcessor(mockService)

	// Create task with invalid payload
	invalidTask := asynq.NewTask(task.TypeEmailNotification, []byte("invalid json"))

	// Act
	err := processor.ProcessTask(context.Background(), invalidTask)

	// Assert
	if err != asynq.SkipRetry {
		t.Errorf("expected SkipRetry, got %v", err)
	}
}

func TestEmailProcessor_ProcessTask_ContextCanceled(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			return ctx.Err()
		},
	}
	processor := NewEmailProcessor(mockService)

	payload := &task.EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	payloadBytes, _ := payload.Marshal()

	asyncTask := asynq.NewTask(task.TypeEmailNotification, payloadBytes)

	// Create canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	err := processor.ProcessTask(ctx, asyncTask)

	// Assert
	if err == nil {
		t.Error("expected context error, got nil")
	}
}

func TestEmailProcessor_ProcessTask_ServiceError(t *testing.T) {
	// Arrange
	callCount := 0
	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			callCount++
			if callCount == 1 {
				return ErrEmailServiceFailed
			}
			return nil
		},
	}
	processor := NewEmailProcessor(mockService)

	payload := &task.EmailNotificationPayload{
		To:      "test@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	payloadBytes, _ := payload.Marshal()

	asyncTask := asynq.NewTask(task.TypeEmailNotification, payloadBytes)

	// Act
	err := processor.ProcessTask(context.Background(), asyncTask)

	// Assert
	if err == nil {
		t.Error("expected service error, got nil")
	}
	if err == asynq.SkipRetry {
		t.Error("expected retryable error, got SkipRetry")
	}
}

func TestEmailProcessor_ProcessTask_ValidatesPayload(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			return nil
		},
	}
	processor := NewEmailProcessor(mockService)

	payload := &task.EmailNotificationPayload{
		To:      "recipient@example.com",
		Subject: "Important Notification",
		Body:    "Please read this important message",
	}
	payloadBytes, _ := payload.Marshal()

	asyncTask := asynq.NewTask(task.TypeEmailNotification, payloadBytes)

	// Act
	err := processor.ProcessTask(context.Background(), asyncTask)

	// Assert
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}
}

func TestEmailProcessor_ProcessTask_EmptyPayload(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{}
	processor := NewEmailProcessor(mockService)

	// Create task with empty payload
	emptyTask := asynq.NewTask(task.TypeEmailNotification, []byte(""))

	// Act
	err := processor.ProcessTask(context.Background(), emptyTask)

	// Assert
	if err != asynq.SkipRetry {
		t.Errorf("expected SkipRetry for empty payload, got %v", err)
	}
}

func TestEmailProcessor_ProcessTask_ParsesPayloadCorrectly(t *testing.T) {
	// Arrange
	receivedTo := ""
	receivedSubject := ""
	receivedBody := ""

	mockService := &MockEmailService{
		SendEmailFunc: func(ctx context.Context, to, subject, body string) error {
			receivedTo = to
			receivedSubject = subject
			receivedBody = body
			return nil
		},
	}
	processor := NewEmailProcessor(mockService)

	expectedPayload := &task.EmailNotificationPayload{
		To:      "user@example.com",
		Subject: "Welcome",
		Body:    "Welcome to our platform",
	}
	payloadBytes, _ := expectedPayload.Marshal()

	asyncTask := asynq.NewTask(task.TypeEmailNotification, payloadBytes)

	// Act
	err := processor.ProcessTask(context.Background(), asyncTask)

	// Assert
	if err != nil {
		t.Fatalf("ProcessTask failed: %v", err)
	}

	if receivedTo != expectedPayload.To {
		t.Errorf("expected to=%q, got %q", expectedPayload.To, receivedTo)
	}
	if receivedSubject != expectedPayload.Subject {
		t.Errorf("expected subject=%q, got %q", expectedPayload.Subject, receivedSubject)
	}
	if receivedBody != expectedPayload.Body {
		t.Errorf("expected body=%q, got %q", expectedPayload.Body, receivedBody)
	}
}

// Helper error for testing
var ErrEmailServiceFailed = &json.SyntaxError{Offset: 0}
