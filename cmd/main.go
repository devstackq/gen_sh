package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devstackq/gen_sh/internal/cron"
	"github.com/devstackq/gen_sh/pkg/logger"
)

func main() {
	// Инициализация логирования
	logger.InitLogger()
	defer logger.Log.Sync()

	// Канал для обработки сигналов (SIGTERM, SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск cron задач для генерации и загрузки видео
	go cron.StartCronJob()

	// Ожидаем сигнала завершения
	sigReceived := <-sigChan
	fmt.Printf("Получен сигнал: %v. Инициируем graceful shutdown...\n", sigReceived)

	// Завершаем работу приложения (позволяя горутинам завершиться)
	// Добавьте здесь любую логику для очистки или завершения работы
	time.Sleep(2 * time.Second)

	fmt.Println("✅ Приложение завершено.")
}
