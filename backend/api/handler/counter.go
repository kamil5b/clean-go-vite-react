package handler

import (
	"net/http"

	"github.com/kamil5b/clean-go-vite-react/backend/service"
	"github.com/labstack/echo/v4"
)

// CounterHandler handles counter-related HTTP requests
type CounterHandler struct {
	service service.CounterService
}

// NewCounterHandler creates a new instance of CounterHandler
func NewCounterHandler(svc service.CounterService) *CounterHandler {
	return &CounterHandler{
		service: svc,
	}
}

// GetCounter handles GET /api/counter requests
func (h *CounterHandler) GetCounter(c echo.Context) error {
	value, err := h.service.GetCounter(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"value": value,
	})
}

// IncrementCounter handles POST /api/counter requests
func (h *CounterHandler) IncrementCounter(c echo.Context) error {
	value, err := h.service.IncrementCounter(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"value": value,
	})
}
