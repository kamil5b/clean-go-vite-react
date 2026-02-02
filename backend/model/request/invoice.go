package request

import "github.com/google/uuid"

type InvoiceItemInput struct {
	ItemID    uuid.UUID `json:"item_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	UnitPrice float64   `json:"unit_price" validate:"required,min=0"`
}

type CreateInvoiceRequest struct {
	GrandPrice float64            `json:"grand_price" validate:"required,min=0"`
	Items      []InvoiceItemInput `json:"items" validate:"required,min=1"`
	Tags       []uuid.UUID        `json:"tags"`
}

type UpdateInvoiceRequest struct {
	GrandPrice float64            `json:"grand_price" validate:"required,min=0"`
	Items      []InvoiceItemInput `json:"items" validate:"required,min=1"`
	Tags       []uuid.UUID        `json:"tags"`
}
