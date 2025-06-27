# Используем официальный образ Go с нужной версией
ARG GO_VERSION=1.24.3
FROM golang:${GO_VERSION}-alpine

# Устанавливаем зависимости для postgres
RUN apk add --no-cache gcc musl-dev

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]