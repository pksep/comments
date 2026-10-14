package api

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type counterRequest struct {
	UserID    string   `json:"user_id" binding:"required"`
	ThreadIDs []string `json:"thread_ids" binding:"required,max=500"`
	MarkRead  bool     `json:"mark_read"`
}

// Counters is an internal endpoint. The ERP gateway supplies the authenticated user ID.
func (h *CommentHandler) Counters(c *gin.Context) {
	secret := os.Getenv("COMMENTS_COUNTERS_KEY")
	if len(secret) < 32 {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "counter authentication unavailable"})
		return
	}
	if subtle.ConstantTimeCompare([]byte(secret), []byte(c.GetHeader("X-Comments-Counters-Key"))) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body counterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid counter request"})
		return
	}
	userID, err := strconv.ParseInt(body.UserID, 10, 64)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user"})
		return
	}
	for _, threadID := range body.ThreadIDs {
		if _, err := uuid.Parse(threadID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread"})
			return
		}
	}

	items, err := h.service.Counters(c.Request.Context(), body.UserID, body.ThreadIDs, body.MarkRead)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "comment counters unavailable"})
		return
	}
	c.JSON(http.StatusOK, items)
}
