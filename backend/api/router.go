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
	notFoundHandler *handler.NotFoundHandler,
	itemHandler *handler.ItemHandler,
	tagHandler *handler.TagHandler,
	invoiceHandler *handler.InvoiceHandler,
) {
	api := e.Group("/api")

	// Public routes (no authentication required)
	api.GET("/message", messageHandler.GetMessage)

	// Auth routes (public)
	api.POST("/auth/register", userHandler.Register)
	api.POST("/auth/login", userHandler.Login)
	api.POST("/auth/refresh", userHandler.Refresh)
	api.GET("/csrf", userHandler.GetCSRFToken)

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(tokenService))

	// Auth protected endpoints
	protected.GET("/auth/me", userHandler.GetMe)

	// Logout requires auth + CSRF protection (it's a POST request)
	protected.POST("/auth/logout", userHandler.Logout, middleware.CSRFMiddleware())

	// Counter endpoints (protected)
	protected.GET("/counter", counterHandler.GetCounter)

	// Counter POST requires auth + CSRF protection
	protected.POST("/counter", counterHandler.IncrementCounter, middleware.CSRFMiddleware())

	// Item endpoints (protected)
	protected.GET("/items", itemHandler.GetAll)
	protected.GET("/items/:id", itemHandler.GetByID)
	protected.POST("/items", itemHandler.Create, middleware.CSRFMiddleware())
	protected.PUT("/items/:id", itemHandler.Update, middleware.CSRFMiddleware())
	protected.DELETE("/items/:id", itemHandler.Delete, middleware.CSRFMiddleware())

	// Tag endpoints (protected)
	protected.GET("/tags", tagHandler.GetAll)
	protected.GET("/tags/:id", tagHandler.GetByID)
	protected.POST("/tags", tagHandler.Create, middleware.CSRFMiddleware())
	protected.PUT("/tags/:id", tagHandler.Update, middleware.CSRFMiddleware())
	protected.DELETE("/tags/:id", tagHandler.Delete, middleware.CSRFMiddleware())

	// Invoice endpoints (protected)
	protected.GET("/invoices", invoiceHandler.GetAll)
	protected.GET("/invoices/:id", invoiceHandler.GetByID)
	protected.POST("/invoices", invoiceHandler.Create, middleware.CSRFMiddleware())
	protected.PUT("/invoices/:id", invoiceHandler.Update, middleware.CSRFMiddleware())
	protected.DELETE("/invoices/:id", invoiceHandler.Delete, middleware.CSRFMiddleware())

	api.Any("/*", notFoundHandler.Handle)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(e *echo.Echo, healthHandler *handler.HealthHandler) {
	api := e.Group("/api")
	api.GET("/health", healthHandler.Check)
}
