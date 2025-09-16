package services

import (
	commentsRepo "github.com/pksep/comments/internal/modules/comments/repository"
	comments "github.com/pksep/comments/internal/modules/comments/service"
)

type Services struct {
	CommentService *comments.CommentService
}

func NewServices(commentRepo commentsRepo.CommentRepoInterface) *Services {
	return &Services{
		CommentService: comments.NewCommentService(commentRepo),
	}
}
