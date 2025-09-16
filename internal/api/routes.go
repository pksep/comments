package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pksep/comments/internal/services"
	commentsApi "github.com/pksep/comments/internal/modules/comments/api"
)

type RouterDeps struct {
}

func RegisterRoutes(r *gin.Engine, deps *RouterDeps, services *services.Services) {
	api := r.Group("/api")

	// Роуты комментариев
	commentHandler := commentsApi.NewCommentHandler(services.CommentService)
	commentHandler.RegisterRoutes(api)
}
