SERVER=cmd/server/main.go
DOCS=internal/api/docs

# --------------------------------------
# Помощь
# --------------------------------------
.PHONY: help
help:
	@echo "Доступные команды Makefile:"
	@echo "  run-dev             - запуск сервера в режиме разработки"
	@echo "  run-prod            - запуск сервера в режиме продакшена"
	@echo "  swag            - генерация Swagger документации"
	@echo "  dev             - генерация Swagger + запуск сервера"
	@echo "  test            - запуск всех тестов"
	@echo "  lint            - запуск golangci-lint"
	@echo "  migration-create - создать новую миграцию (NAME=description)"

# --------------------------------------
# Сборка и запуск сервера
# --------------------------------------
.PHONY: run-dev
run-dev:
	air -c .air.toml

.PHONY: run-prod
run-prod:
	go run $(SERVER)

# --------------------------------------
# Генерация Swagger документации
# --------------------------------------
.PHONY: swag
swag:
	swag init -g $(SERVER) -o $(DOCS)

# --------------------------------------
# Генерация + запуск
# --------------------------------------
.PHONY: dev
dev: swag run

# --------------------------------------
# Тесты
# --------------------------------------
.PHONY: test
test:
	go test ./...

# --------------------------------------
# Линтеры
# --------------------------------------
.PHONY: lint
lint:
	golangci-lint run

# --------------------------------------
# Миграции базы данных
# --------------------------------------
MIGRATIONS_DIR = migrations

.PHONY: migration-create
migration-create:
	@if [ -z "${NAME}" ]; then \
		echo "Error: NAME is not set. Usage: make migration-create NAME=description_of_changes"; \
		exit 1; \
	fi; \
	TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq "$${TIMESTAMP}_${NAME}" -digits 14
