package tag

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
)

// GetByID gets a tag by ID
func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*response.TagResponse, error) {
	tag, err := s.tagRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &response.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		ColorHex:  tag.ColorHex,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, nil
}

// GetAll gets all tags with pagination
func (s *tagService) GetAll(ctx context.Context, page, limit int, search string) (*response.TagPaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	tags, total, err := s.tagRepository.FindAll(ctx, page, limit, search)
	if err != nil {
		return nil, err
	}

	tagResponses := make([]response.TagResponse, len(tags))
	for i, tag := range tags {
		tagResponses[i] = response.TagResponse{
			ID:        tag.ID,
			Name:      tag.Name,
			ColorHex:  tag.ColorHex,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		}
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return &response.TagPaginationResponse{
		Data: tagResponses,
		Meta: response.TagPaginationMeta{
			TotalData: int(total),
			Page:      page,
			Limit:     limit,
			TotalPage: totalPage,
		},
	}, nil
}
