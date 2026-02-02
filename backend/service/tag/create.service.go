package tag

import (
	"context"
	"errors"

	"github.com/kamil5b/clean-go-vite-react/backend/model/entity"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// Create creates a new tag
func (s *tagService) Create(ctx context.Context, req *request.CreateTagRequest) (*response.TagResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.ColorHex == "" {
		return nil, errors.New("color_hex is required")
	}

	tag := entity.TagEntity{
		Name:     req.Name,
		ColorHex: req.ColorHex,
	}

	id, err := s.tagRepository.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	created, err := s.tagRepository.FindByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return &response.TagResponse{
		ID:        created.ID,
		Name:      created.Name,
		ColorHex:  created.ColorHex,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}
