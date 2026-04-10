package api

import (
	"net/http"
	"strings"

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
		comments.POST("/create", h.Create)
		comments.POST("/update", h.Update) // id будет в теле
		comments.POST("/delete", h.Delete) // id, author_id будет в теле
		comments.POST("/pin", h.Pin)
		comments.POST("/unpin", h.Unpin)
		comments.GET("/by-thread/:threadId", h.Get)
		comments.GET("/list", h.List) // ids[]=id1&ids[]=id2
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var body dto.CreateCommentDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(body.Content) == "" && len(body.Documents) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content or documents is required"})
		return
	}

	created, err := h.service.Create(c, model.Comment{
		AuthorID:        body.AuthorID,
		Content:         body.Content,
		ThreadID:        body.ThreadID,
		AnswerCommentID: body.AnswerCommentID,
		Documents: func() []model.CommentMedia {
			documents := make([]model.CommentMedia, 0, len(body.Documents))
			for _, document := range body.Documents {
				documents = append(documents, model.CommentMedia{
					ID:           document.ID,
					Name:         document.Name,
					OriginalName: document.OriginalName,
					Type:         document.Type,
					Size:         document.Size,
				})
			}

			return documents
		}(),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *CommentHandler) Update(c *gin.Context) {
	var body dto.UpdateCommentDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.UpdateContent(c, body.ID, body.Content, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	var body dto.DeleteCommentDTO

	// Проверяем входные данные
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Удаляем комментарий
	deletedComment, err := h.service.Delete(c, body.ID, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем удалённый комментарий
	c.JSON(http.StatusOK, deletedComment)
}

func (h *CommentHandler) Pin(c *gin.Context) {
	var body dto.PinCommentDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.Pin(c, body.ID, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *CommentHandler) Unpin(c *gin.Context) {
	var body dto.PinCommentDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.Unpin(c, body.ID, body.AuthorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *CommentHandler) Get(c *gin.Context) {
	threadId := c.Param("threadId")
	item, err := h.service.GetByID(c, threadId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *CommentHandler) List(c *gin.Context) {
	idsParam := c.Query("ids")
	var ids []string
	if idsParam != "" {
		ids = strings.Split(idsParam, ",")
	}

	// тут limit = 3 для реплаев
	items, err := h.service.ListWithReplies(c, ids, 3)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}
