package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pksep/location_search_server/internal/modules/users/api/dto"
	users "github.com/pksep/location_search_server/internal/modules/users/service"
)

type UserHandler struct {
	service *users.UserService
}

func NewUserHandler(service *users.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.ListUsers)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

// CreateUser создает нового пользователя
// @Summary      Создать пользователя
// @Description  Создаёт нового пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserDTO  true  "Пользователь"
// @Success      201   {object}  model.User
// @Failure      400   {object}  map[string]interface{} "ошибка валидации"
// @Failure      500   {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var body dto.CreateUserDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.service.CreateUser(c, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// ListUsers выводит список пользователей
// @Summary      Список пользователей
// @Description  Возвращает всех пользователей
// @Tags         users
// @Produce      json
// @Success      200  {array}   model.User
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {


	users, err := h.service.ListUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser возвращает одного пользователя по ID
// @Summary      Получить пользователя
// @Description  Возвращает пользователя по ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object}  model.User
// @Failure      404  {object}  map[string]interface{} "пользователь не найден"
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser обновляет данные пользователя
// @Summary      Обновить пользователя
// @Description  Обновляет инициалы пользователя по ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "ID пользователя"
// @Param        user  body      dto.UpdateUserDTO true  "Пользователь"
// @Success      200   {object}  model.User
// @Failure      400   {object}  map[string]interface{} "ошибка валидации"
// @Failure      500   {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var body dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	updatedUser, err := h.service.UpdateUser(c, id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}


// DeleteUser удаляет пользователя
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по ID
// @Tags         users
// @Param        id   path      string  true  "ID пользователя"
// @Success      204  "No Content"
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
