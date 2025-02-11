# Базовый образ
FROM golang:1.19-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы в контейнер
COPY . .

# Устанавливаем зависимости
RUN go mod tidy

# Компиляция приложения
RUN go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем cron
RUN apk update && apk add --no-cache cron

# Копируем скомпилированное приложение и конфигурационные файлы
COPY --from=builder /app/main /app/main
COPY config.yaml /app/config.yaml

# Копируем crontab файл
COPY cronfile /etc/crontabs/root

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем точку входа
CMD ["sh", "-c", "crond -f -L /dev/stdout & /app/main"]
