package config

import (
	"log"
	"os"
	"strings"
	"sync"
)

type Config struct {
	DatabaseURL        string
	Port               string
	MinioPublicBaseURL string
	MinioBucketName    string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			log.Fatal("DATABASE_URL не задан")
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		minioPublicBaseURL := strings.TrimRight(os.Getenv("MINIO_PUBLIC_BASE_URL"), "/")
		if minioPublicBaseURL == "" {
			log.Fatal("MINIO_PUBLIC_BASE_URL не задан")
		}

		minioBucketName := strings.Trim(os.Getenv("MINIO_BUCKET_NAME"), "/")
		if minioBucketName == "" {
			log.Fatal("MINIO_BUCKET_NAME не задан")
		}

		instance = &Config{
			DatabaseURL:        dbURL,
			Port:               port,
			MinioPublicBaseURL: minioPublicBaseURL,
			MinioBucketName:    minioBucketName,
		}
	})
	return instance
}
