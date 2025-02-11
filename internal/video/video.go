package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/upload"
	"github.com/devstackq/gen_sh/pkg/logger"
	"github.com/pkg/errors"
)

// GenerateAndPublishForUser генерирует и публикует видео для пользователя на его платформах
func GenerateAndPublishForUser(user config.User) error {

	audioPath, err := GenerateAudioForText(user.Theme)
	if err != nil {
		log.Printf("Ошибка генерации аудио для темы %s: %v", user.Theme, err)
		return err
	}

	// Генерация видео
	videoPath, err := GenerateVideo(user.Theme, audioPath)
	if err != nil {
		log.Printf("Ошибка генерации видео для пользователя %s: %v", user.Email, err)
		return err
	}

	// Публикация видео на платформы
	var wg sync.WaitGroup
	for _, platform := range user.Platforms {
		wg.Add(1)
		go func(platform config.Platform) {
			defer wg.Done()
			err := upload.UploadToPlatform(platform, videoPath)
			if err != nil {
				log.Printf("Ошибка при публикации видео для пользователя %s на платформе %s: %v", user.Email, platform.Name, err)
			}
		}(platform)
	}

	wg.Wait()

	return nil
}

// GenerateVideo - основная логика генерации видео на основе текста и аудио.
func GenerateVideo(text, audioPath string) (string, error) {
	// 1. Генерация временного видео с текстом.
	textVideoPath, err := generateTextVideo(text)
	if err != nil {
		return "", errors.Wrap(err, "ошибка генерации видео с текстом")
	}

	// 2. Наложение аудио на видео.
	videoPath, err := combineAudioWithVideo(textVideoPath, audioPath)
	if err != nil {
		return "", errors.Wrap(err, "ошибка наложения аудио на видео")
	}

	// 3. Очистка временных файлов.
	err = removeFile(textVideoPath)
	if err != nil {
		logger.Log.Warn("Не удалось удалить временное видео с текстом", "file", textVideoPath)
	}

	// Возвращаем путь к финальному видео.
	return videoPath, nil
}

// generateTextVideo - генерирует видео с наложением текста.
func generateTextVideo(text string) (string, error) {
	// Создаем временный путь для видео.
	videoPath := fmt.Sprintf("/tmp/%d_text_video.mp4", time.Now().Unix())

	// Команда для создания простого видео с текстом.
	// Используем ffmpeg или другую подходящую утилиту.
	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-t", "5", "-i", "color=c=black:s=1280x720:r=30", "-vf", fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=48:x=(w-text_w)/2:y=(h-text_h)/2", text), "-an", videoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("не удалось создать видео с текстом: %v, output: %s", err, string(output))
	}

	return videoPath, nil
}

func GenerateAudioForText(text string) (string, error) {
	// Генерация аудио для текста
	// Здесь мы вызываем реальную логику генерации аудио (например, с помощью Google TTS или другого сервиса)
	audioPath, err := speech.GenerateAudio(text)
	if err != nil {
		log.Printf("Ошибка генерации аудио для текста %s: %v", text, err)
		return "", err
	}
	return audioPath, nil
}

// combineAudioWithVideo - объединяет аудиофайл и видеофайл в одно.
func combineAudioWithVideo(videoPath, audioPath string) (string, error) {
	// Создаем путь для итогового видео.
	finalVideoPath := fmt.Sprintf("/tmp/%d_final_video.mp4", time.Now().Unix())

	// Команда для объединения видео с аудио.
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-shortest", finalVideoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("не удалось объединить видео с аудио: %v, output: %s", err, string(output))
	}

	// Возвращаем путь к итоговому видео.
	return finalVideoPath, nil
}

// removeFile удаляет файл и логирует ошибку, если она возникла.
func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("ошибка при удалении файла %s: %v", path, err)
	}
	return nil
}
