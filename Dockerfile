# --- Этап сборки ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum (чтобы кешировать зависимости)
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o server main.go

# --- Этап запуска ---
FROM alpine:latest
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/server .

# Порт API
EXPOSE 8080

CMD ["./server"]
