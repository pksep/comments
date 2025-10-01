package db

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pksep/comments/internal/config"
)

// ----------------------------------------------------------------------
// RunMigrations																				  							|
// ищет папку migrations относительно internal/db и запускает миграции  |
// ----------------------------------------------------------------------
func RunMigrations() {
	cfg := config.GetConfig()
	dbURL := cfg.DatabaseURL

	// Определяем путь к текущему файлу (migrate.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Не удалось определить путь к текущему файлу")
	}

	// migrations расположена в корне проекта
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationsPath := filepath.Join(projectRoot, "migrations")

	log.Printf("Применяем миграции из: %s\n", migrationsPath)

	// Normalize to file URL (cross-OS). On Windows this becomes file:///C:/...
	slashed := filepath.ToSlash(migrationsPath)
	var sourceURL string
	if runtime.GOOS == "windows" {
		sourceURL = "file://" + strings.TrimPrefix(slashed, "/")
	} else {
		sourceURL = "file:///" + strings.TrimPrefix(slashed, "/")
	}

	m, err := migrate.New(
		sourceURL,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Ошибка при создании миграции: %v", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Миграции уже применены, изменений нет")
		} else {
			log.Fatalf("Ошибка при применении миграций: %v", err)
		}
	} else {
		log.Println("Миграции успешно применены")
	}
}
