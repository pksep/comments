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
		comments.POST("/create", h.Create)
		comments.POST("/update", h.Update)   // id будет в теле
		comments.POST("/delete", h.Delete)   // id будет в теле
		comments.GET("/:id", h.Get)
		comments.GET("/list", h.List)        // ids[]=id1&ids[]=id2
	}
}

func (h *CommentHandler) Create(c *gin.Context) {
	var body dto.CreateCommentDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.Create(c, model.Comment{
		AuthorID:        body.AuthorID,
		Content:         body.Content,
		ParentCommentID: body.ParentCommentID,
		AnswerCommentID:  body.AnswerCommentID,
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

	updated, err := h.service.UpdateContent(c, body.ID, body.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	var body struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c, body.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CommentHandler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *CommentHandler) List(c *gin.Context) {
    ids := c.QueryArray("ids[]")

    // тут limit = 3 для реплаев
    items, err := h.service.ListWithReplies(c, ids, 3)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, items)
}

