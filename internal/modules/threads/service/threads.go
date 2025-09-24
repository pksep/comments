package service

import (
	"github.com/pksep/comments/internal/modules/threads/repository"
)

type ThreadService struct {
	repo repository.ThreadRepoInterface
}

func NewThreadService(repo repository.ThreadRepoInterface) *ThreadService {
	return &ThreadService{repo: repo}
}