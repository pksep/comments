package config

import (
	"github.com/gin-gonic/gin"
	"github.com/pksep/location_search_server/internal/ws"
)

func RegisterRoutes(r *gin.Engine) {
	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)

	r.GET("/ws", wsHandler.Connect)
}