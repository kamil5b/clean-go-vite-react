package handler

import (
	"net/http"

	"github.com/kamil5b/clean-go-vite-react/internal/service"
	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	service service.HealthService
}

// NewHealthHandler creates a new instance of HealthHandler
func NewHealthHandler(svc service.HealthService) *HealthHandler {
	return &HealthHandler{
		service: svc,
	}
}

// Check handles GET /api/health requests
func (h *HealthHandler) Check(c echo.Context) error {
	status, err := h.service.Check(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, status)
}

// CheckWithDependencies handles GET /api/health/detailed requests
func (h *HealthHandler) CheckWithDependencies(c echo.Context, checks map[string]func(c echo.Context) error) error {
	checkFuncs := make(map[string]func(c echo.Context) error)
	for name, check := range checks {
		checkFuncs[name] = check
	}

	status, err := h.service.Check(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, status)
}
