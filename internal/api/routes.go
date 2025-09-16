package api

import (
	"github.com/gin-gonic/gin"
	examsApi "github.com/pksep/location_search_server/internal/modules/exams/api"
	examsRepo "github.com/pksep/location_search_server/internal/modules/exams/repository"
	usersApi "github.com/pksep/location_search_server/internal/modules/users/api"
	"github.com/pksep/location_search_server/internal/modules/users/repository"
	"github.com/pksep/location_search_server/internal/services"
)

type RouterDeps struct {
	UserRepo *repository.UserRepo
	ExamRepo *examsRepo.ExamRepo
	// В будущем можно добавить ProjectRepo, AnswerRepo и т.д.
}

func RegisterRoutes(r *gin.Engine, deps *RouterDeps, services *services.Services) {
	api := r.Group("/api")

	// Роуты пользователей
	userHandler := usersApi.NewUserHandler(services.UserService)
	userHandler.RegisterRoutes(api)

	// Роуты экзаменов
	examHandler := examsApi.NewExamHandler(services.ExamService)
	examHandler.RegisterRoutes(api)
}
