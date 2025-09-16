package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/pksep/comments/internal/modules/comments/api/dto"
    "github.com/pksep/comments/internal/modules/comments/model"
    comments "github.com/pksep/comments/internal/modules/comments/service"
)

type CommentHandler struct {
    service *comments.CommentService
}

func NewCommentHandler(service *comments.CommentService) *CommentHandler {
    return &CommentHandler{service: service}
}

func (h *CommentHandler) RegisterRoutes(rg *gin.RouterGroup) {
    comments := rg.Group("/comments")
    {
        comments.POST("", h.Create)
        comments.GET("", h.ListByEntity)
        comments.GET(":id", h.Get)
        comments.PUT(":id", h.Update)
        comments.DELETE(":id", h.Delete)
    }
}

func (h *CommentHandler) Create(c *gin.Context) {
    // @Summary      Создать комментарий
    // @Description  Создаёт комментарий к сущности
    // @Tags         comments
    // @Accept       json
    // @Produce      json
    // @Param        comment  body      dto.CreateCommentDTO  true  "Комментарий"
    // @Success      201      {object}  model.Comment
    // @Failure      400      {object}  map[string]interface{} "ошибка валидации"
    // @Router       /comments [post]
    var body dto.CreateCommentDTO
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    created, err := h.service.Create(c, model.Comment{
        EntityType: body.EntityType,
        EntityID:   body.EntityID,
        AuthorID:   body.AuthorID,
        Content:    body.Content,
        ParentID:   body.ParentID,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
}

func (h *CommentHandler) ListByEntity(c *gin.Context) {
    // @Summary      Список комментариев по сущности
    // @Description  Возвращает комментарии по паре entity_type + entity_id
    // @Tags         comments
    // @Produce      json
    // @Param        entity_type  query  string  true  "Тип сущности"
    // @Param        entity_id    query  string  true  "ID сущности"
    // @Success      200  {array}   model.Comment
    // @Failure      400  {object}  map[string]interface{} "ошибка валидации"
    // @Router       /comments [get]
    entityType := c.Query("entity_type")
    entityID := c.Query("entity_id")
    items, err := h.service.ListByEntity(c, entityType, entityID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, items)
}

func (h *CommentHandler) Get(c *gin.Context) {
    // @Summary      Получить комментарий
    // @Description  Возвращает комментарий по ID
    // @Tags         comments
    // @Produce      json
    // @Param        id   path      string  true  "ID комментария"
    // @Success      200  {object}  model.Comment
    // @Failure      404  {object}  map[string]interface{} "не найдено"
    // @Router       /comments/{id} [get]
    id := c.Param("id")
    item, err := h.service.GetByID(c, id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, item)
}

func (h *CommentHandler) Update(c *gin.Context) {
    // @Summary      Обновить комментарий
    // @Description  Обновляет текст комментария по ID
    // @Tags         comments
    // @Accept       json
    // @Produce      json
    // @Param        id       path      string            true  "ID комментария"
    // @Param        comment  body      dto.UpdateCommentDTO  true  "Комментарий"
    // @Success      200      {object}  model.Comment
    // @Failure      400      {object}  map[string]interface{} "ошибка валидации"
    // @Router       /comments/{id} [put]
    id := c.Param("id")
    var body dto.UpdateCommentDTO
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    updated, err := h.service.UpdateContent(c, id, body.Content)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, updated)
}

func (h *CommentHandler) Delete(c *gin.Context) {
    // @Summary      Удалить комментарий
    // @Description  Удаляет комментарий по ID
    // @Tags         comments
    // @Param        id   path      string  true  "ID комментария"
    // @Success      204  "No Content"
    // @Failure      500  {object}  map[string]interface{} "внутренняя ошибка сервера"
    // @Router       /comments/{id} [delete]
    id := c.Param("id")
    if err := h.service.Delete(c, id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}


