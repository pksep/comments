package repository

import (
	"context"

	"github.com/pksep/location_search_server/internal/modules/slides/model"
)

type SlidesRepoInterface interface {
	Create(ctx context.Context, slide *model.Slide) (*model.Slide, error)
	GetList(ctx context.Context) ([]model.Slide, error)
	GetByID(ctx context.Context, id string) (*model.Slide, error)
	Update(ctx context.Context, slide *model.Slide) (*model.Slide, error)
	Delete(ctx context.Context, id string) error
}