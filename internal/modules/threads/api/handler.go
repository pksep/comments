package api

import (
	"github.com/pksep/comments/internal/modules/threads/service"
)

type ThreadHandler struct {
	service *service.ThreadService
}

func NewThreadHandler(service *service.ThreadService) *ThreadHandler {
	return &ThreadHandler{service: service}
}
