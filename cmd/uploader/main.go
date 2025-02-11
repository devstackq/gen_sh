package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/devstackq/gen_sh/internal/database"
	"github.com/devstackq/gen_sh/internal/upload"
	"github.com/devstackq/gen_sh/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	if err := database.InitDB(); err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	fmt.Println("🚀 Запуск загрузчика видео...")

	rows, err := database.DB.Query("SELECT id, path FROM videos WHERE uploaded = false")
	if err != nil {
		log.Fatal("Ошибка получения видео из БД:", err)
	}
	defer rows.Close()

	var wg sync.WaitGroup

	for rows.Next() {
		var id int
		var videoPath string
		err := rows.Scan(&id, &videoPath)
		if err != nil {
			log.Println("Ошибка чтения строки:", err)
			continue
		}

		wg.Add(1)
		go func(id int, videoPath string) {
			defer wg.Done()
			err := upload.UploadToYouTube(videoPath)
			if err == nil {
				database.DB.Exec("UPDATE videos SET uploaded = true WHERE id = $1", id)
			}
		}(id, videoPath)
	}

	wg.Wait()
	fmt.Println("✅ Все видео загружены!")
}
