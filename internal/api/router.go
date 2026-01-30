package api

import (
	"github.com/kamil5b/clean-go-vite-react/internal/api/handler"
	"github.com/kamil5b/clean-go-vite-react/internal/service"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, messageService service.MessageService) {
	api := e.Group("/api")

	// Health check endpoint
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Message endpoints
	messageHandler := handler.NewMessageHandler(messageService)
	api.GET("/message", messageHandler.GetMessage)
}
