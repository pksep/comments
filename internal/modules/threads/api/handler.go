package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pksep/comments/internal/modules/comments/model"
	"github.com/pksep/comments/internal/modules/comments/service"
)

type ThreadHandler struct {
	service *service.ThreadService
}

func NewThreadHandler(service *service.ThreadService) *ThreadHandler {
	return &ThreadHandler{service: service}
}

func (h *ThreadHandler) RegisterRoutes(rg *gin.RouterGroup) {
	threads := rg.Group("/threads")
	{
		threads.POST("/create", h.Create)
	}
}

func (h *ThreadHandler) Create(c *gin.Context) {
	thread, err := h.service.Create(c, model.Thread{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, thread)
}
