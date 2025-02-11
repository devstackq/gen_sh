package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/devstackq/gen_sh/internal/content"
	"github.com/devstackq/gen_sh/internal/database"
	"github.com/devstackq/gen_sh/internal/monitoring"
	"github.com/devstackq/gen_sh/internal/speech"
	"github.com/devstackq/gen_sh/internal/textprocessing"
	"github.com/devstackq/gen_sh/internal/video"
	"github.com/devstackq/gen_sh/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	if err := database.InitDB(); err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	monitoring.InitMonitoring()
	fmt.Println("🚀 Запуск генератора видео...")

	text, err := content.FetchNews()
	if err != nil {
		log.Fatalf("Ошибка парсинга контента: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	var processedText string
	go func() {
		defer wg.Done()
		processedText = textprocessing.ProcessText(text)
	}()

	wg.Add(1)
	var audioPath string
	go func() {
		defer wg.Done()
		var err error
		audioPath, err = speech.GenerateAudio(processedText)
		if err != nil {
			log.Printf("Ошибка генерации аудио: %v", err)
		}
	}()

	wg.Wait()

	videoPath, err := video.CreateVideo(processedText, audioPath)
	if err != nil {
		log.Fatalf("Ошибка создания видео: %v", err)
	}

	fmt.Printf("✅ Видео создано: %s\n", videoPath)

	_, err = database.DB.Exec("INSERT INTO videos (title, path, created_at) VALUES ($1, $2, $3)",
		processedText, videoPath, time.Now())
	if err != nil {
		log.Printf("Ошибка сохранения видео в БД: %v", err)
	}

	monitoring.VideosGenerated.Inc()
	logger.Log.Info("🎬 Видео успешно сохранено в БД")
}
