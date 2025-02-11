# Используем официальный образ Go для сборки
FROM golang:1.19 AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь проект
COPY . .

# Собираем проект
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/gen_sh cmd/generator/main.go

# Используем минимальный образ для исполнения
FROM debian:bullseye-slim

# Устанавливаем необходимые утилиты (например, для Cron)
RUN apt-get update && apt-get install -y cron

# Копируем скомпилированный бинарник из стадии сборки
COPY --from=build /bin/gen_sh /bin/gen_sh

# Копируем shutdown скрипт для graceful shutdown
COPY shutdown.sh /usr/local/bin/shutdown.sh
RUN chmod +x /usr/local/bin/shutdown.sh

# Устанавливаем точку входа
ENTRYPOINT ["/bin/gen_sh"]
