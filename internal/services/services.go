package services

import (
	commentsRepo "github.com/pksep/comments/internal/modules/comments/repository"
	threadsRepo "github.com/pksep/comments/internal/modules/threads/repository"

	commentsSvc "github.com/pksep/comments/internal/modules/comments/service"
	threadsSvc  "github.com/pksep/comments/internal/modules/threads/service"
)

// Services объединяет все бизнес-сервисы
type Services struct {
	CommentService *commentsSvc.CommentService
	ThreadService  *threadsSvc.ThreadService
}

// NewServices конструктор, принимает репозитории и возвращает набор сервисов
func NewServices(
	commentRepo commentsRepo.CommentRepoInterface,
	threadRepo threadsRepo.ThreadRepoInterface,
) *Services {
	return &Services{
		CommentService: commentsSvc.NewCommentService(commentRepo),
		ThreadService:  threadsSvc.NewThreadService(threadRepo),
	}
}
