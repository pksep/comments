package service

import (
    "context"
    "errors"

    "github.com/pksep/comments/internal/modules/comments/model"
    "github.com/pksep/comments/internal/modules/comments/repository"
)

type CommentService struct {
    repo repository.CommentRepoInterface
}

func NewCommentService(repo repository.CommentRepoInterface) *CommentService {
    return &CommentService{repo: repo}
}

func (s *CommentService) Create(ctx context.Context, input model.Comment) (*model.Comment, error) {
    if input.EntityType == "" || input.EntityID == "" {
        return nil, errors.New("entity reference is required")
    }
    if input.AuthorID == "" {
        return nil, errors.New("author_id is required")
    }
    if input.Content == "" {
        return nil, errors.New("content is required")
    }
    return s.repo.Create(ctx, &input)
}

func (s *CommentService) GetByID(ctx context.Context, id string) (*model.Comment, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *CommentService) UpdateContent(ctx context.Context, id string, content string) (*model.Comment, error) {
    if content == "" {
        return nil, errors.New("content is required")
    }
    existing, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    existing.Content = content
    return s.repo.Update(ctx, existing)
}

func (s *CommentService) Delete(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

func (s *CommentService) ListByEntity(ctx context.Context, entityType string, entityID string) ([]model.Comment, error) {
    if entityType == "" || entityID == "" {
        return nil, errors.New("entity reference is required")
    }
    return s.repo.ListByEntity(ctx, entityType, entityID)
}


