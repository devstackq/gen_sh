package cron

import (
	"fmt"
	"log"
	"sync"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/upload"
	"github.com/devstackq/gen_sh/internal/video"
	"github.com/robfig/cron/v3"
)

// StartCronJob запускает cron задачи для генерации и публикации видео
func StartCronJob(cfg *config.Config) {
	// Создание нового cron экземпляра
	c := cron.New(cron.WithSeconds())

	// Добавляем задачу для генерации видео каждый день в 1 час ночи
	_, err := c.AddFunc("0 1 * * *", func() {
		fmt.Println("🚀 Запуск задачи по генерации видео...")

		// Генерация и публикация для каждого пользователя
		var wg sync.WaitGroup
		for _, user := range cfg.Users {
			wg.Add(1)
			go func(user config.User) {
				defer wg.Done()
				err := video.GenerateAndPublishForUser(user)
				if err != nil {
					log.Printf("Ошибка при генерации и публикации видео для пользователя %s: %v", user.Email, err)
				}
			}(user)
		}

		wg.Wait()
	})

	// Добавляем задачу для публикации видео каждый день в 2 часа ночи
	_, err = c.AddFunc("0 2 * * *", func() {
		fmt.Println("🚀 Запуск задачи по загрузке видео...")

		// Публикация видео на платформы для каждого пользователя
		var wg sync.WaitGroup
		for _, user := range cfg.Users {
			wg.Add(1)
			go func(user config.User) {
				defer wg.Done()
				err := video.GenerateAndPublishForUser(user)
				if err != nil {
					log.Printf("Ошибка при генерации и публикации видео для пользователя %s: %v", user.Email, err)
				}
			}(user)
		}

		wg.Wait()
	})

	// Запуск cron задач
	c.Start()

	// Ожидание завершения всех задач
	select {}
}
