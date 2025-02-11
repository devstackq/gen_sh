# Базовый образ Go
FROM golang:1.20 AS builder

# Установка необходимых системных пакетов (espeak, ffmpeg, python3, pip)
RUN apt-get update && apt-get install -y \
    espeak \
    ffmpeg \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Установка gTTS (Google Text-to-Speech)
RUN pip3 install gtts

# Создание рабочей директории
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Сборка бинарного файла
RUN go mod tidy && go build -o main ./cmd/main.go

# Финальный образ (минимальный, для продакшна)
FROM debian:bullseye-slim

# Установка зависимостей для ffmpeg и espeak (только нужные библиотеки)
RUN apt-get update && apt-get install -y \
    espeak \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Копируем собранное приложение из builder
WORKDIR /app
COPY --from=builder /app/main .

# Запуск приложения
CMD [ "./main" ]
