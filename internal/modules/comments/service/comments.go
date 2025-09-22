package service

import (
	"context"

	"github.com/pksep/comments/internal/modules/comments/model"
	"github.com/pksep/comments/internal/modules/comments/repository"
)

type CommentService struct {
	repo repository.CommentRepoInterface
}

func NewCommentService(repo repository.CommentRepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(ctx context.Context, c model.Comment) (*model.Comment, error) {
	return s.repo.Create(ctx, &c)
}

func (s *CommentService) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) UpdateContent(ctx context.Context, id string, content string) (*model.Comment, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.Content = content
	return s.repo.Update(ctx, c)
}

func (s *CommentService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CommentService) ListWithReplies(ctx context.Context, ids []string, depth int) ([]model.Comment, error) {
	return s.repo.ListWithReplies(ctx, ids, depth)
}
