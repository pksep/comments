package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pksep/location_search_server/internal/api"
	"github.com/pksep/location_search_server/internal/config"
	examRepoPkg "github.com/pksep/location_search_server/internal/modules/exams/repository"
	userRepoPkg "github.com/pksep/location_search_server/internal/modules/users/repository"
	"github.com/pksep/location_search_server/internal/services"
)

func Init(pool *pgxpool.Pool) *gin.Engine {

	// Инициализация репозиториев
	userRepo := userRepoPkg.NewUserRepo(pool)
	examRepo := examRepoPkg.NewExamRepo(pool)
	// projectRepo := repository.NewProjectRepo(pool) // если будет

	// Инициализация сервисов
	services := services.NewServices(userRepo, examRepo)

	// Инициализация зависимостей для хэндлеров
	deps := &api.RouterDeps{
		UserRepo: userRepo,
		ExamRepo: examRepo,
		// ProjectRepo: projectRepo,
	}

	// Инициализация Gin
	r := gin.Default()

	// Swagger
	swaggerCfg := &config.SwaggerConfig{
		Title:       "Exam Management API",
		Description: "API для управления экзаменами",
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
