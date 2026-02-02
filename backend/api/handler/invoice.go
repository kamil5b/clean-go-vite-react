package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	invoiceSvc "github.com/kamil5b/clean-go-vite-react/backend/service/invoice"
	"github.com/labstack/echo/v4"
)

// InvoiceHandler handles invoice-related HTTP requests
type InvoiceHandler struct {
	invoiceService invoiceSvc.InvoiceService
}

// NewInvoiceHandler creates a new instance of InvoiceHandler
func NewInvoiceHandler(invoiceService invoiceSvc.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

// Create handles POST /api/invoices requests
func (h *InvoiceHandler) Create(c echo.Context) error {
	req := &request.CreateInvoiceRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	invoice, err := h.invoiceService.Create(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, invoice)
}

// GetByID handles GET /api/invoices/:id requests
func (h *InvoiceHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	invoice, err := h.invoiceService.GetByID(c.Request().Context(), uuid.UUID(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, invoice)
}

// Update handles PUT /api/invoices/:id requests
func (h *InvoiceHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	req := &request.UpdateInvoiceRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	invoice, err := h.invoiceService.Update(c.Request().Context(), uuid.UUID(id), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, invoice)
}

// Delete handles DELETE /api/invoices/:id requests
func (h *InvoiceHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	if err := h.invoiceService.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "invoice deleted successfully",
	})
}

// GetAll handles GET /api/invoices requests
func (h *InvoiceHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	search := c.QueryParam("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	invoices, err := h.invoiceService.GetAll(c.Request().Context(), page, limit, search)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, invoices)
}
