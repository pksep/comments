package config

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/pksep/comments/internal/api/docs"
)

// SwaggerConfig хранит основные параметры Swagger
type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
}

// GetSwaggerHandler возвращает Gin handler для Swagger UI
func GetSwaggerHandler(cfg *SwaggerConfig) gin.HandlerFunc {
	// Настройка документации через swag
	if cfg.Title != "" {
		docs.SwaggerInfo.Title = cfg.Title
	}
	if cfg.Description != "" {
		docs.SwaggerInfo.Description = cfg.Description
	}
	if cfg.Version != "" {
		docs.SwaggerInfo.Version = cfg.Version
	}
	if cfg.Host != "" {
		docs.SwaggerInfo.Host = cfg.Host
	}
	if cfg.BasePath != "" {
		docs.SwaggerInfo.BasePath = cfg.BasePath
	}
	if len(cfg.Schemes) > 0 {
		docs.SwaggerInfo.Schemes = cfg.Schemes
	}

	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
