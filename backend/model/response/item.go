package response

import (
	"time"

	"github.com/google/uuid"
)

type ItemResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ItemPaginationMeta struct {
	TotalData int `json:"totalData"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"totalPage"`
}

type ItemPaginationResponse struct {
	Data []ItemResponse     `json:"data"`
	Meta ItemPaginationMeta `json:"meta"`
}
