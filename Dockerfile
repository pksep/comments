FROM golang:1.25-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download
ARG ENV_FILE=.env
# Копируем исходники
COPY . .
COPY ${ENV_FILE} ./.env
# Устанавливаем air для hot reload (если нужно в dev)
RUN go install github.com/air-verse/air@latest

# Собираем бинарь
RUN go build -o server ./cmd/server

# Минимальный рантайм-образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарь из builder
COPY --from=builder /app/server .
COPY --from=builder /app/migrations /app/migrations

# (опционально) если нужен air в dev:
COPY --from=builder /go/bin/air /usr/local/bin/air
