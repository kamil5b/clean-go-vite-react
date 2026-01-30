package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kamil5b/clean-go-vite-react/internal/api/handler"
	"github.com/kamil5b/clean-go-vite-react/internal/service"
	"github.com/labstack/echo/v4"
)

func TestSetupRoutes_HealthEndpoint(t *testing.T) {
	// Arrange
	e := echo.New()
	messageService := service.NewMessageService()
	healthService := service.NewHealthService()
	SetupRoutes(e, messageService)

	// Register health handler (done in DI container normally)
	healthHandler := handler.NewHealthHandler(healthService)
	e.GET("/api/health", healthHandler.Check)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", response["status"])
	}
}

func TestSetupRoutes_MessageEndpoint(t *testing.T) {
	// Arrange
	e := echo.New()
	messageService := service.NewMessageService()
	SetupRoutes(e, messageService)

	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "Hello, from the golang World!"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q, got %q", expectedMessage, response["message"])
	}
}

func TestSetupRoutes_NotFound(t *testing.T) {
	// Arrange
	e := echo.New()
	messageService := service.NewMessageService()
	SetupRoutes(e, messageService)

	req := httptest.NewRequest(http.MethodGet, "/api/nonexistent", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

type MockMessageServiceForIntegration struct {
	message string
	err     error
}

func (m *MockMessageServiceForIntegration) GetMessage(ctx context.Context) (string, error) {
	return m.message, m.err
}

func TestSetupRoutes_WithMockedService(t *testing.T) {
	// Arrange
	e := echo.New()
	mockService := &MockMessageServiceForIntegration{
		message: "Mocked Message",
		err:     nil,
	}
	SetupRoutes(e, mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
	rec := httptest.NewRecorder()

	// Act
	e.ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if response["message"] != "Mocked Message" {
		t.Errorf("expected 'Mocked Message', got %q", response["message"])
	}
}

func TestSetupRoutes_MultipleRequests(t *testing.T) {
	// Arrange
	e := echo.New()
	messageService := service.NewMessageService()
	SetupRoutes(e, messageService)

	// Act & Assert - make multiple requests
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/message", nil)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d", i, http.StatusOK, rec.Code)
		}

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("request %d: failed to unmarshal response: %v", i, err)
		}

		if response["message"] != "Hello, from the golang World!" {
			t.Errorf("request %d: unexpected message", i)
		}
	}
}
