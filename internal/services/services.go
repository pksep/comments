package services

import (
	usersRepo "github.com/pksep/location_search_server/internal/modules/users/repository"
	users "github.com/pksep/location_search_server/internal/modules/users/service"

	examsRepo "github.com/pksep/location_search_server/internal/modules/exams/repository"
	exams "github.com/pksep/location_search_server/internal/modules/exams/service"
)

type Services struct {
	UserService *users.UserService
	ExamService *exams.ExamService
}

func NewServices(userRepo usersRepo.UserRepoInterface, examRepo examsRepo.ExamRepoInterface) *Services {
	return &Services{
		UserService: users.NewUserService(userRepo),
		ExamService: exams.NewExamService(examRepo),
	}
}
