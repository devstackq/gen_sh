package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/cron"
	"github.com/devstackq/gen_sh/pkg/logger"
)

func main() {
	// Инициализация логирования
	logger.InitLogger()
	defer logger.Log.Sync()

	// Загрузка конфигурации из файла
	configFile := "config.yaml"
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Запуск cron задач
	go cron.StartCronJob(cfg)

	// Ожидаем завершения всех cron задач
	select {}
}
