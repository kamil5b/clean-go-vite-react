package response

import (
	"time"

	"github.com/google/uuid"
)

type InvoiceItemResponse struct {
	ID         uuid.UUID    `json:"id"`
	ItemID     uuid.UUID    `json:"item_id"`
	Item       ItemResponse `json:"item"`
	Quantity   int          `json:"quantity"`
	UnitPrice  float64      `json:"unit_price"`
	TotalPrice float64      `json:"total_price"`
}

type InvoiceResponse struct {
	ID         uuid.UUID `json:"id"`
	GrandPrice float64   `json:"grand_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type InvoiceDetailResponse struct {
	ID         uuid.UUID             `json:"id"`
	GrandPrice float64               `json:"grand_price"`
	Items      []InvoiceItemResponse `json:"items"`
	Tags       []TagResponse         `json:"tags"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type InvoiceListItem struct {
	ID         uuid.UUID     `json:"id"`
	GrandPrice float64       `json:"grand_price"`
	Tags       []TagResponse `json:"tags"`
	TotalItem  int           `json:"totalItem"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

type InvoicePaginationMeta struct {
	TotalData int `json:"totalData"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
}

type InvoicePaginationResponse struct {
	Data []InvoiceListItem     `json:"data"`
	Meta InvoicePaginationMeta `json:"meta"`
}
