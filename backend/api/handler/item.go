package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	itemSvc "github.com/kamil5b/clean-go-vite-react/backend/service/item"
	"github.com/labstack/echo/v4"
)

// ItemHandler handles item-related HTTP requests
type ItemHandler struct {
	itemService itemSvc.ItemService
}

// NewItemHandler creates a new instance of ItemHandler
func NewItemHandler(itemService itemSvc.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// Create handles POST /api/items requests
func (h *ItemHandler) Create(c echo.Context) error {
	req := &request.CreateItemRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	item, err := h.itemService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, item)
}

// GetByID handles GET /api/items/:id requests
func (h *ItemHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	item, err := h.itemService.GetByID(c.Request().Context(), uuid.UUID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, item)
}

// Update handles PUT /api/items/:id requests
func (h *ItemHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	req := &request.UpdateItemRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	item, err := h.itemService.Update(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, item)
}

// Delete handles DELETE /api/items/:id requests
func (h *ItemHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	if err := h.itemService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "item deleted successfully",
	})
}

// GetAll handles GET /api/items requests
func (h *ItemHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items, err := h.itemService.GetAll(c.Request().Context(), page, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, items)
}
