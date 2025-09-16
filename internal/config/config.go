package config

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	DatabaseURL string
	Port        string
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

		instance = &Config{
			DatabaseURL: dbURL,
			Port:        port,
		}
	})
	return instance
}
