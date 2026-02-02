package api

import (
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/kamil5b/clean-go-vite-react/backend/api/middleware"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(e *echo.Echo, logic handler.Logic) {
	// Initialize handler
	healthHandler := handler.NewHealthHandler()
	messageHandler := handler.NewMessageHandler(logic)
	counterHandler := handler.NewCounterHandler(logic)
	userHandler := handler.NewUserHandler(logic)
	notFoundHandler := handler.NewNotFoundHandler()

	// API group
	api := e.Group("/api")

	// Health check
	api.GET("/health", healthHandler.Check)

	// Public routes
	api.GET("/message", messageHandler.GetMessage)
	api.GET("/counter", counterHandler.GetCounter)
	api.POST("/counter", counterHandler.IncrementCounter)

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/logout", userHandler.Logout)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(logic))
	protected.GET("/auth/me", userHandler.GetMe)

	// 404 handler for undefined API endpoints (must be last)
	api.Any("/*", notFoundHandler.Handle)
}
