package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/api"
	"github.com/pksep/comments/internal/config"
	commentRepoPkg "github.com/pksep/comments/internal/modules/comments/repository"
	threadRepoPkg "github.com/pksep/comments/internal/modules/threads/repository"
	"github.com/pksep/comments/internal/services"
)

func Init(pool *pgxpool.Pool) *gin.Engine {

	// Инициализация репозиториев
	commentRepo := commentRepoPkg.NewCommentRepo(pool)
	threadRepo := threadRepoPkg.NewThreadRepo(pool)

	// Инициализация сервисов
	services := services.NewServices(commentRepo, threadRepo)

	// Инициализация зависимостей для хэндлеров
	deps := &api.RouterDeps{}

	// Инициализация Gin
	r := gin.Default()

	// Swagger
	swaggerCfg := &config.SwaggerConfig{
		Title:       "Comments API",
		Description: "Универсальный сервис комментариев",
		Version:     "1.0.0",
		Host:        "localhost:5001",
		BasePath:    "/api",
		Schemes:     []string{"http"},
	}
	r.GET("/swagger/*any", config.GetSwaggerHandler(swaggerCfg))

	// ...existing code...

	// Регистрация всех API маршрутов
	api.RegisterRoutes(r, deps, services)

	return r
}
