package tag

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamil5b/clean-go-vite-react/backend/model/request"
	"github.com/kamil5b/clean-go-vite-react/backend/model/response"
	"github.com/kamil5b/clean-go-vite-react/backend/repository/interfaces"
)

// TagService defines the interface for tag operations
type TagService interface {
	Create(ctx context.Context, req *request.CreateTagRequest) (*response.TagResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*response.TagResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *request.UpdateTagRequest) (*response.TagResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, page, limit int, search string) (*response.TagPaginationResponse, error)
}

// tagService is the concrete implementation of TagService
type tagService struct {
	tagRepository interfaces.TagRepository
}

// NewTagService creates a new instance of TagService
func NewTagService(tagRepository interfaces.TagRepository) TagService {
	return &tagService{
		tagRepository: tagRepository,
	}
}
