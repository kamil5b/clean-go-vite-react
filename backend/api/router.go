package api

import (
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	e *echo.Echo,
	messageHandler handler.MessageHandler,
	counterHandler handler.CounterHandler,
) {
	api := e.Group("/api")

	api.GET("/message", messageHandler.GetMessage)

	api.GET("/counter", counterHandler.GetCounter)
	api.POST("/counter", counterHandler.IncrementCounter)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(e *echo.Echo, healthHandler *handler.HealthHandler) {
	api := e.Group("/api")
	api.GET("/health", healthHandler.Check)
}
