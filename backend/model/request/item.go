package request

type CreateItemRequest struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
}

type UpdateItemRequest struct {
	Name string `json:"name" validate:"required"`
	Desc string `json:"desc"`
}
