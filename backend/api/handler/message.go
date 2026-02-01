package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	logic Logic
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(logic Logic) *MessageHandler {
	return &MessageHandler{logic: logic}
}

// GetMessage handles GET /api/message
func (h *MessageHandler) GetMessage(c echo.Context) error {
	content, err := h.logic.GetMessage(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"content": content,
	})
}
