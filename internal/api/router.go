package api

import (
	"github.com/kamil5b/clean-go-vite-react/internal/api/handler"
	"github.com/kamil5b/clean-go-vite-react/internal/service"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, messageService service.MessageService, counterService service.CounterService) {
	api := e.Group("/api")

	// Message endpoints
	messageHandler := handler.NewMessageHandler(messageService)
	api.GET("/message", messageHandler.GetMessage)

	// Counter endpoints
	counterHandler := handler.NewCounterHandler(counterService)
	api.GET("/counter", counterHandler.GetCounter)
	api.POST("/counter", counterHandler.IncrementCounter)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(e *echo.Echo, healthHandler *handler.HealthHandler) {
	api := e.Group("/api")
	api.GET("/health", healthHandler.Check)
}
