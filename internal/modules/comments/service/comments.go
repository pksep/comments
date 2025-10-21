package service

import (
	"context"

	"github.com/pksep/comments/internal/modules/comments/model"
	"github.com/pksep/comments/internal/modules/comments/repository"
)

type CommentService struct {
	repo repository.CommentRepoInterface
}

// NewCommentService создаёт новый сервис комментариев
func NewCommentService(repo repository.CommentRepoInterface) *CommentService {
	return &CommentService{repo: repo}
}

// Create создаёт новый комментарий
func (s *CommentService) Create(ctx context.Context, c model.Comment) (*model.Comment, error) {
	return s.repo.Create(ctx, &c)
}

// GetByID возвращает комментарий по threadId
func (s *CommentService) GetByID(ctx context.Context, threadId string) (*model.Comment, error) {
	return s.repo.GetByID(ctx, threadId)
}

// UpdateContent обновляет контент комментария
func (s *CommentService) UpdateContent(ctx context.Context, id string, content string, authorId string) (*model.Comment, error) {
	return s.repo.Update(ctx, id, content, authorId)
}

// Delete удаляет комментарий
func (s *CommentService) Delete(ctx context.Context, id string, authorId string) (*model.Comment, error) {
	return s.repo.Delete(ctx, id, authorId)
}

// ListWithReplies возвращает root-комменты с ограничением replyLimit реплаев
func (s *CommentService) ListWithReplies(ctx context.Context, ids []string, replyLimit int) ([]model.Comment, error) {
	return s.repo.ListWithReplies(ctx, ids, replyLimit)
}
