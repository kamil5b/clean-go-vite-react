package handler

import (
	"net/http"

	"github.com/kamil5b/clean-go-vite-react/backend/service"
	"github.com/labstack/echo/v4"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	service service.MessageService
}

// NewMessageHandler creates a new instance of MessageHandler
func NewMessageHandler(svc service.MessageService) *MessageHandler {
	return &MessageHandler{
		service: svc,
	}
}

// GetMessage handles GET /api/message requests
func (h *MessageHandler) GetMessage(c echo.Context) error {
	message, err := h.service.GetMessage(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})
}
