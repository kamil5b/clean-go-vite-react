package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CounterHandler handles counter-related HTTP requests
type CounterHandler struct {
	logic Logic
}

// NewCounterHandler creates a new counter handler
func NewCounterHandler(logic Logic) *CounterHandler {
	return &CounterHandler{logic: logic}
}

// GetCounter handles GET /api/counter
func (h *CounterHandler) GetCounter(c echo.Context) error {
	value, err := h.logic.GetCounter(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"value": value,
	})
}

// IncrementCounter handles POST /api/counter
func (h *CounterHandler) IncrementCounter(c echo.Context) error {
	value, err := h.logic.IncrementCounter(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"value": value,
	})
}
