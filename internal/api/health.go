package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	dbPool *pgxpool.Pool
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func NewHealthHandler(dbPool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		dbPool: dbPool,
	}
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response := HealthResponse{
		Timestamp: time.Now(),
		Services:  make(map[string]string),
	}

	allHealthy := true

	// Check database connectivity
	if h.dbPool != nil {
		if err := h.dbPool.Ping(ctx); err != nil {
			response.Services["database"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			response.Services["database"] = "healthy"
		}
	} else {
		response.Services["database"] = "not initialized"
		allHealthy = false
	}

	if allHealthy {
		response.Status = "ready"
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = "not ready"
		c.JSON(http.StatusServiceUnavailable, response)
	}
}
