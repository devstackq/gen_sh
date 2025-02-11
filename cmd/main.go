package main

import (
	"log"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/cron"
	"github.com/devstackq/gen_sh/internal/logger"
)

func main() {

	if err := logger.InitLogger("logs/app.log"); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.CloseLogger()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	go cron.StartCronJob(cfg)

	// Ожидаем завершения всех cron задач
	select {}
}
