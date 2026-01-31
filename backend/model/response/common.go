package response

import "github.com/google/uuid"

type CommonIDResponse struct {
	ID uuid.UUID `json:"value"`
}
