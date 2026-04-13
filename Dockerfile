# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

ARG ENV_FILE=.env
COPY . .
# Копируем нужный конфиг под стандартным именем
COPY ${ENV_FILE} ./.env

# Собираем статический бинарник (CGO_ENABLED=0 важен для alpine)
# Убедитесь, что точка входа действительно в ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/comments_server ./cmd/server

# --- Stage 2: Runtime ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Забираем только результат сборки и конфиги
COPY --from=builder /app/comments_server .
COPY --from=builder /app/.env ./.env

# Если у комментариев есть свои миграции
COPY --from=builder /app/migrations ./migrations
