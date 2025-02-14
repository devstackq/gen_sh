package cron

import (
	"fmt"
	"log"
	"sync"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/content"
	"github.com/devstackq/gen_sh/internal/video"
)

func Start(cfg *config.Config) {
	// Создание нового cron экземпляра
	//c := cron.New(cron.WithSeconds())

	// Добавляем задачу для генерации видео каждый день в 1 час ночи
	//_, err := c.AddFunc("0 1 0 * * *", func() { // todo - flexible cron - RUN CROND

	fmt.Println("🚀 Запуск задачи по генерации видео...")

	// Генерация и публикация для каждого пользователя
	var wg sync.WaitGroup
	for _, user := range cfg.Users {

		wg.Add(1)
		go func(user config.User) {

			defer wg.Done()
			items, err := content.FetchContent(user.Theme, user.Sources)
			if err != nil {
				log.Fatalf("Ошибка получения контента: %v", err)
			}
			if len(items) == 0 {
				log.Fatalf("content is empty")
			}

			videoPath, err := video.GenerateVideo(user, items)
			if err != nil {
				log.Fatalf("GenerateVideo %s: %v", user.Email, err)
			}
			items[0].Path = videoPath

			if err = video.Publish(user, items[0]); err != nil {
				log.Fatalf("Publish %s: %v", user.Email, err)
			}

		}(user)
	}

	wg.Wait()
	//})
	//if err != nil {
	//	fmt.Println("Cron error.", err)
	//}

	// Запуск cron задач
	//c.Start()

	// Ожидание завершения всех задач
	select {}
}
