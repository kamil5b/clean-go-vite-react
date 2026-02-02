package response

import (
	"time"

	"github.com/google/uuid"
)

type TagResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ColorHex  string    `json:"color_hex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagPaginationMeta struct {
	TotalData int `json:"totalData"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
}

type TagPaginationResponse struct {
	Data []TagResponse     `json:"data"`
	Meta TagPaginationMeta `json:"meta"`
}
