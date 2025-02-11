package video

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/devstackq/gen_sh/internal/video/model"
	"github.com/devstackq/gen_sh/pkg/logger"
)

// CreateVideo - основная логика генерации видео на основе текста и аудио.
func CreateVideo(text, audioPath string) (string, error) {
	// 1. Генерация временного видео с текстом.
	textVideoPath, err := generateTextVideo(text)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации видео с текстом: %v", err)
	}

	// 2. Наложение аудио на видео.
	videoPath, err := combineAudioWithVideo(textVideoPath, audioPath)
	if err != nil {
		return "", fmt.Errorf("ошибка наложения аудио на видео: %v", err)
	}

	// 3. Очистка временных файлов.
	err = os.Remove(textVideoPath)
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

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("не удалось создать видео с текстом: %v", err)
	}

	return videoPath, nil
}

// combineAudioWithVideo - объединяет аудиофайл и видеофайл в одно.
func combineAudioWithVideo(videoPath, audioPath string) (string, error) {
	// Создаем путь для итогового видео.
	finalVideoPath := fmt.Sprintf("/tmp/%d_final_video.mp4", time.Now().Unix())

	// Команда для объединения видео с аудио.
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-shortest", finalVideoPath)

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("не удалось объединить видео с аудио: %v", err)
	}

	// Возвращаем путь к итоговому видео.
	return finalVideoPath, nil
}
