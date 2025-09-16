package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pksep/comments/internal/app"
	"github.com/pksep/comments/internal/config"
	"github.com/pksep/comments/internal/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Нет .env файла, используем системные переменные")
	}
}

func main() {
	cfg := config.GetConfig()

	// Создаём пул подключений к Postgres
	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	// Автоматический запуск миграций
	db.RunMigrations()

	// Инициализация Gin
	r := app.Init(pool)	

	// Запуск сервера
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

}
