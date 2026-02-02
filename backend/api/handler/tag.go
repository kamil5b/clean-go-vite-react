package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	tagSvc "github.com/kamil5b/clean-go-vite-react/backend/service/tag"
	"github.com/labstack/echo/v4"
)

// TagHandler handles tag-related HTTP requests
type TagHandler struct {
	tagService tagSvc.TagService
}

// NewTagHandler creates a new instance of TagHandler
func NewTagHandler(tagService tagSvc.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// Create handles POST /api/tags requests
func (h *TagHandler) Create(c echo.Context) error {
	req := &request.CreateTagRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	tag, err := h.tagService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, tag)
}

// GetByID handles GET /api/tags/:id requests
func (h *TagHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	tag, err := h.tagService.GetByID(c.Request().Context(), uuid.UUID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tag)
}

// Update handles PUT /api/tags/:id requests
func (h *TagHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	req := &request.UpdateTagRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	tag, err := h.tagService.Update(c.Request().Context(), uuid.UUID(id), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tag)
}

// Delete handles DELETE /api/tags/:id requests
func (h *TagHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	if err := h.tagService.Delete(c.Request().Context(), uuid.UUID(id)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "tag deleted successfully",
	})
}

// GetAll handles GET /api/tags requests
func (h *TagHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	tags, err := h.tagService.GetAll(c.Request().Context(), page, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, tags)
}
