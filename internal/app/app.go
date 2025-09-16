package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/comments/internal/api"
	"github.com/pksep/comments/internal/config"
	"github.com/pksep/comments/internal/services"
	commentRepoPkg "github.com/pksep/comments/internal/modules/comments/repository"
)

func Init(pool *pgxpool.Pool) *gin.Engine {

	// Инициализация репозиториев
	commentRepo := commentRepoPkg.NewCommentRepo(pool)

	// Инициализация сервисов
	services := services.NewServices(commentRepo)

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

	// Подключаем WS-модуль
	config.RegisterRoutes(r)

	// Регистрация всех API маршрутов
	api.RegisterRoutes(r, deps, services)

	return r
}
