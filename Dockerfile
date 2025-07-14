# Шаг 1: Используем официальный образ Go для сборки
FROM golang:1.24.3-alpine AS builder
# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum (если есть)
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
# -ldflags="-s -w" - убираем отладочную информацию для уменьшения размера бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Шаг 2: Используем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранный бинарник из builder-образа
COPY --from=builder /app/main .

# Копируем конфигурационные файлы
COPY configs/ /root/configs/

# Копируем .env файл
COPY .env .

# Команда для запуска приложения
CMD ["./main"]