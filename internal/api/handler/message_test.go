package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// MockMessageService is a mock implementation of MessageService
type MockMessageService struct {
	GetMessageFunc func(ctx context.Context) (string, error)
}

func (m *MockMessageService) GetMessage(ctx context.Context) (string, error) {
	if m.GetMessageFunc != nil {
		return m.GetMessageFunc(ctx)
	}
	return "", nil
}

func TestMessageHandler_GetMessage_Success(t *testing.T) {
	// Arrange
	mockService := &MockMessageService{
		GetMessageFunc: func(ctx context.Context) (string, error) {
			return "Hello, from the golang World!", nil
		},
	}
	handler := NewMessageHandler(mockService)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := handler.GetMessage(c)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expectedBody := `{"message":"Hello, from the golang World!"}`
	if rec.Body.String() != expectedBody+"\n" {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestMessageHandler_GetMessage_ServiceError(t *testing.T) {
	// Arrange
	mockService := &MockMessageService{
		GetMessageFunc: func(ctx context.Context) (string, error) {
			return "", errors.New("service error")
		},
	}
	handler := NewMessageHandler(mockService)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := handler.GetMessage(c)

	// Assert
	if err != nil {
		t.Fatalf("expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	if !contains(rec.Body.String(), "service error") {
		t.Errorf("expected error message in body, got %s", rec.Body.String())
	}
}

func TestMessageHandler_GetMessage_ContextCanceled(t *testing.T) {
	// Arrange
	mockService := &MockMessageService{
		GetMessageFunc: func(ctx context.Context) (string, error) {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			default:
				return "test", nil
			}
		},
	}
	handler := NewMessageHandler(mockService)

	e := echo.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := handler.GetMessage(c)

	// Assert
	if err != nil {
		t.Fatalf("expected no error from handler, got %v", err)
	}

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
