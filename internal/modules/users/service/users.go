package users

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/pksep/location_search_server/internal/modules/users/api/dto"
	"github.com/pksep/location_search_server/internal/modules/users/model"
	"github.com/pksep/location_search_server/internal/modules/users/repository"
)

type UserService struct {
	repo repository.UserRepoInterface
}


func NewUserService(repo repository.UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx *gin.Context,body dto.CreateUserDTO) (model.User, error) {

	// TODO проверить уникальность инициалов

	user := &model.User{
		Initials: body.Initials,
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return model.User{}, errors.New("Ошибка:" + err.Error())
	}

	return *createdUser, nil
}

func (s *UserService) ListUsers(ctx *gin.Context) ([]model.User, error) {
	users, err := s.repo.List(ctx.Request.Context())
	if err != nil {
		return []model.User{}, errors.New("Ошибка:" + err.Error())
	}

	return  users, nil
}

func (s *UserService) GetUser(ctx *gin.Context, id string) (model.User, error) {

	user, err := s.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		return model.User{}, errors.New("Ошибка:" + err.Error())
		
	}

	return *user, nil
}

func (s *UserService) UpdateUser(ctx *gin.Context, id string, body dto.UpdateUserDTO) (model.User, error) {
	user := &model.User{
		ID:       id,
		Initials: body.Initials,
	}

	updatedUser, err := s.repo.Update(ctx.Request.Context(), user)
	if err != nil {
		return model.User{}, errors.New("Ошибка:" + err.Error())
	}

	return *updatedUser, nil
}

func (s *UserService) DeleteUser(ctx *gin.Context, id string) error {
	err := s.repo.Delete(ctx.Request.Context(), id)
	if err != nil {
		return errors.New("Ошибка:" + err.Error())
	}
	return nil
}