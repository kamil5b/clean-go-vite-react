package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ItemHandler handles item-related HTTP requests
type ItemHandler struct {
	logic Logic
}

// NewItemHandler creates a new item handler
func NewItemHandler(logic Logic) *ItemHandler {
	return &ItemHandler{logic: logic}
}

// CreateItem handles POST /api/items
func (h *ItemHandler) CreateItem(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
	}

	item, err := h.logic.CreateItem(c.Request().Context(), req.Title, req.Description, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, ItemToInfo(item))
}

// GetItems handles GET /api/items
func (h *ItemHandler) GetItems(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	items, err := h.logic.GetItems(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ItemsToInfoList(items))
}

// GetItem handles GET /api/items/:id
func (h *ItemHandler) GetItem(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid item id",
		})
	}

	item, err := h.logic.GetItemByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
	}

	// Check ownership
	if item.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "forbidden",
		})
	}

	return c.JSON(http.StatusOK, ItemToInfo(item))
}

// UpdateItem handles PUT /api/items/:id
func (h *ItemHandler) UpdateItem(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid item id",
		})
	}

	// Check ownership first
	existingItem, err := h.logic.GetItemByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
	}

	if existingItem.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "forbidden",
		})
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "title is required",
		})
	}

	item, err := h.logic.UpdateItem(c.Request().Context(), id, req.Title, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ItemToInfo(item))
}

// DeleteItem handles DELETE /api/items/:id
func (h *ItemHandler) DeleteItem(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid item id",
		})
	}

	// Check ownership
	item, err := h.logic.GetItemByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "item not found",
		})
	}

	if item.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "forbidden",
		})
	}

	if err := h.logic.DeleteItem(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
