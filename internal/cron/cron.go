package cron

import (
	"fmt"
	"log"

	"github.com/devstackq/gen_sh/internal/database"
	"github.com/devstackq/gen_sh/internal/upload"
	"github.com/devstackq/gen_sh/internal/video/service"
	"github.com/devstackq/gen_sh/pkg/logger"
	"github.com/robfig/cron/v3"
)

// StartCronJob запускает задачи cron для поиска, генерации и загрузки видео по расписанию
func StartCronJob() {
	// Создание нового cron экземпляра
	c := cron.New(cron.WithSeconds())

	// Запуск задачи для генерации видео каждый день в 1 час ночи (пример)
	_, err := c.AddFunc("0 1 * * *", func() {
		fmt.Println("🚀 Запуск задачи по генерации видео...")
		// Генерация видео
		if err := GenerateVideos(); err != nil {
			log.Printf("Ошибка при генерации видео по cron: %v", err)
		}
	})

	// Запуск задачи для загрузки видео каждый день в 2 часа ночи (пример)
	_, err = c.AddFunc("0 2 * * *", func() {
		fmt.Println("🚀 Запуск задачи по загрузке видео...")
		// Публикация видео
		if err := PublishVideos(); err != nil {
			log.Printf("Ошибка при публикации видео по cron: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Ошибка добавления cron задачи: %v", err)
	}

	// Запуск cron задач
	c.Start()
}

// GenerateVideos генерирует новые видео
func GenerateVideos() error {
	// Инициализация подключения к базе данных
	dbConn, err := database.ConnectDB()
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return err
	}

	// Создание сервиса видео
	videoService := service.NewVideoService(dbConn)

	// Получаем контент для видео (например, новости, истории, факты)
	// Это будет зависеть от логики вашего контента
	content, err := service.FetchContent()
	if err != nil {
		log.Printf("Ошибка получения контента: %v", err)
		return err
	}

	// Генерация видео на основе контента
	for _, c := range content {
		videoPath, err := service.GenerateVideo(c)
		if err != nil {
			log.Printf("Ошибка при генерации видео для контента %v: %v", c, err)
			continue
		}

		// Сохраняем сгенерированное видео в базе данных
		err = videoService.SaveGeneratedVideo(c, videoPath)
		if err != nil {
			log.Printf("Ошибка при сохранении видео в базе данных: %v", err)
		}
	}

	return nil
}

// PublishVideos выполняет загрузку видео на платформу (например, YouTube)
func PublishVideos() error {
	// Инициализация подключения к базе данных
	dbConn, err := database.ConnectDB()
	if err != nil {
		log.Printf("Ошибка подключения к БД: %v", err)
		return err
	}

	// Создание сервиса видео
	videoService := service.NewVideoService(dbConn)

	// Получаем видео, которые нужно загрузить на платформу
	videos, err := videoService.GetUnpublishedVideos()
	if err != nil {
		log.Printf("Ошибка получения видео из базы данных: %v", err)
		return err
	}

	// Инициализация загрузчика YouTube
	credentialsFile := "credentials.json"
	ytUploader, err := upload.NewYouTubeUploader(credentialsFile)
	if err != nil {
		log.Printf("Ошибка инициализации YouTube uploader: %v", err)
		return err
	}

	// Параллельная загрузка видео на YouTube
	for _, videoRecord := range videos {
		go func(videoRecord service.Video) {
			// Загрузка видео на платформу
			err := ytUploader.Upload(videoRecord.Path)
			if err != nil {
				log.Printf("Ошибка при загрузке видео %s: %v", videoRecord.Path, err)
				return
			}

			// Обновление статуса видео в базе данных
			err = videoService.MarkVideoAsPublished(videoRecord.ID)
			if err != nil {
				log.Printf("Ошибка обновления статуса видео в БД: %v", err)
			} else {
				log.Printf("Видео %s успешно загружено на YouTube", videoRecord.Path)
			}
		}(videoRecord)
	}

	return nil
}
