package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pksep/location_search_server/internal/modules/exams/api/dto"
	exams "github.com/pksep/location_search_server/internal/modules/exams/service"
)

type ExamHandler struct {
	service *exams.ExamService
}

func NewExamHandler(service *exams.ExamService) *ExamHandler {
	return &ExamHandler{service: service}
}

func (h *ExamHandler) RegisterRoutes(rg *gin.RouterGroup) {
	exams := rg.Group("/exams")
	{
		exams.POST("", h.CreateExam)
		exams.GET("", h.ListExams)
		exams.GET("/:id", h.GetExam)
		exams.PUT("/:id", h.UpdateExam)
		exams.DELETE("/:id", h.DeleteExam)
	}
}

// CreateExam создает новый экзамен
// @Summary      Создать экзамен
// @Description  Создает новый экзамен
// @Tags         exams
// @Accept       json
// @Produce      json
// @Param        exam  body      dto.CreateExamDTO  true  "Экзамен"
// @Success      201   {object}  model.Exam
// @Failure      400   {object}  map[string]interface{} "ошибка валидации"
// @Failure      500   {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /exams [post]
func (h *ExamHandler) CreateExam(c *gin.Context) {
	var body dto.CreateExamDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdExam, err := h.service.CreateExam(c, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdExam)
}

// ListExams выводит список экзаменов
// @Summary      Список экзаменов
// @Description  Возвращает все экзамены
// @Tags         exams
// @Produce      json
// @Success      200  {array}   model.Exam
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /exams [get]
func (h *ExamHandler) ListExams(c *gin.Context) {
	exams, err := h.service.ListExams(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exams)
}

// GetExam возвращает один экзамен по ID
// @Summary      Получить экзамен
// @Description  Возвращает экзамен по ID
// @Tags         exams
// @Produce      json
// @Param        id   path      string  true  "ID экзамена"
// @Success      200  {object}  model.Exam
// @Failure      404  {object}  map[string]interface{} "экзамен не найден"
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /exams/{id} [get]
func (h *ExamHandler) GetExam(c *gin.Context) {
	id := c.Param("id")
	exam, err := h.service.GetExam(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exam)
}

// UpdateExam обновляет данные экзамена
// @Summary      Обновить экзамен
// @Description  Обновляет данные экзамена по ID
// @Tags         exams
// @Accept       json
// @Produce      json
// @Param        id    path      string           true  "ID экзамена"
// @Param        exam  body      dto.UpdateExamDTO true  "Экзамен"
// @Success      200   {object}  model.Exam
// @Failure      400   {object}  map[string]interface{} "ошибка валидации"
// @Failure      500   {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /exams/{id} [put]
func (h *ExamHandler) UpdateExam(c *gin.Context) {
	id := c.Param("id")

	var body dto.UpdateExamDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedExam, err := h.service.UpdateExam(c, id, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedExam)
}

// DeleteExam удаляет экзамен
// @Summary      Удалить экзамен
// @Description  Удаляет экзамен по ID
// @Tags         exams
// @Param        id   path      string  true  "ID экзамена"
// @Success      204  "No Content"
// @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
// @Router       /exams/{id} [delete]
func (h *ExamHandler) DeleteExam(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteExam(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
