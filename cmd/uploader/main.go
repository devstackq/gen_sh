package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/devstackq/gen_sh/internal/database"
	"github.com/devstackq/gen_sh/internal/upload"
	"github.com/devstackq/gen_sh/internal/video/service"
	"github.com/devstackq/gen_sh/pkg/logger"
)

func main() {
	// Инициализация логирования
	logger.InitLogger()
	defer logger.Log.Sync()

	// Путь к credentials для YouTube
	credentialsFile := "credentials.json"

	// Инициализация загрузчика YouTube
	ytUploader, err := upload.NewYouTubeUploader(credentialsFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации YouTube uploader: %v", err)
	}
	//todo future upload tikTok

	// Подключение к базе данных
	dbConn, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Создание сервиса для видео
	videoService := service.NewVideoService(dbConn)

	// Получение видео, которые не были загружены
	videos, err := videoService.GetUnuploadedVideos()
	if err != nil {
		log.Fatal("Ошибка получения видео из БД:", err)
	}

	// Запуск параллельной загрузки видео на YouTube
	fmt.Println("🚀 Запуск загрузчика видео...")

	var wg sync.WaitGroup
	for _, videoRecord := range videos {
		wg.Add(1)
		go func(videoRecord service.Video) {
			defer wg.Done()

			// Загрузка видео на YouTube
			err := ytUploader.Upload(videoRecord.Path)
			if err != nil {
				log.Printf("Ошибка загрузки видео %s: %v", videoRecord.Path, err)
				return
			}

			// Обновление статуса загрузки видео в базе данных
			err = videoService.MarkVideoAsUploaded(videoRecord.ID)
			if err != nil {
				log.Printf("Ошибка обновления статуса видео в БД: %v", err)
			} else {
				log.Printf("Видео %s успешно загружено на YouTube", videoRecord.Path)
			}
		}(videoRecord)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	fmt.Println("✅ Все видео загружены!")
}
