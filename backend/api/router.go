package api

import (
	"github.com/kamil5b/clean-go-vite-react/backend/api/handler"
	"github.com/kamil5b/clean-go-vite-react/backend/api/middleware"
	"github.com/kamil5b/clean-go-vite-react/backend/service/token"
	"github.com/labstack/echo/v4"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	e *echo.Echo,
	messageHandler handler.MessageHandler,
	counterHandler handler.CounterHandler,
	userHandler *handler.UserHandler,
	tokenService token.TokenService,
) {
	api := e.Group("/api")

	// Public routes
	api.GET("/message", messageHandler.GetMessage)
	api.GET("/counter", counterHandler.GetCounter)
	api.POST("/counter", counterHandler.IncrementCounter)

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.GET("/auth/csrf", userHandler.GetCSRFToken)
	api.POST("/auth/logout", userHandler.Logout)
	api.POST("/auth/refresh", userHandler.Refresh)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))
	protected.GET("/auth/me", userHandler.GetMe)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(e *echo.Echo, healthHandler *handler.HealthHandler) {
	api := e.Group("/api")
	api.GET("/health", healthHandler.Check)
}
