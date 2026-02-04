package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/services"
	commentsApi "github.com/pksep/comments/internal/modules/comments/api"
)

type RouterDeps struct {
}

func RegisterRoutes(r *gin.Engine, deps *RouterDeps, services *services.Services, dbPool *pgxpool.Pool) {
	// Health check endpoint
	healthHandler := NewHealthHandler(dbPool)
	r.GET("/ready", healthHandler.Ready)

	api := r.Group("/api")

	// Роуты комментариев
	commentHandler := commentsApi.NewCommentHandler(services.CommentService)
	commentHandler.RegisterRoutes(api)
}
