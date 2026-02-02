package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// NotFoundHandler handles 404 responses for API endpoints
type NotFoundHandler struct{}

// NewNotFoundHandler creates a new not found handler
func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

// Handle handles GET /api/* for undefined endpoints
func (h *NotFoundHandler) Handle(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"error":   "Not Found",
		"message": "The requested API endpoint does not exist",
		"path":    c.Request().URL.Path,
		"method":  c.Request().Method,
	})
}
