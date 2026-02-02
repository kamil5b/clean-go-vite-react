package request

type CreateTagRequest struct {
	Name     string `json:"name" validate:"required"`
	ColorHex string `json:"color_hex" validate:"required"`
}

type UpdateTagRequest struct {
	Name     string `json:"name" validate:"required"`
	ColorHex string `json:"color_hex" validate:"required"`
}
