package video

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/content"
	"github.com/devstackq/gen_sh/internal/logger"
	"github.com/devstackq/gen_sh/internal/speech"
	"github.com/devstackq/gen_sh/internal/uploader"
	"github.com/pkg/errors"
)

type Video struct {
	uploader uploader.PlatformClient
}

func Publish(user config.User, content []content.Content) error {

	logger.LogInfo(fmt.Sprint("Начата обработка пользователя", "email", user.Email, "theme", user.Theme))
	// Генерация аудио
	audioPath, err := GenerateAudioForText(content[0].Excerpt) //todo mb use 1 content?
	if err != nil {
		return err
	}

	// Генерация видео
	videoPath, err := GenerateVideo(content[0].Excerpt, audioPath)
	if err != nil {
		return err
	}

	//title, description, tags := GenerateMetadata(user.Theme)

	// Инициализация клиентов для всех платформ пользователя
	clients := make(map[string]uploader.PlatformClient)
	for _, platform := range user.Platforms {
		client, err := uploader.New(platform)
		if err != nil {
			logger.LogError(fmt.Sprintf("Ошибка инициализации клиента %s: %v", platform.Name, err))
			continue
		}
		clients[platform.Name] = client
	}

	// Параллельная публикация на платформы
	var wg sync.WaitGroup
	for _, platform := range user.Platforms {
		wg.Add(1)
		go func(platform config.Platform) {
			defer wg.Done()
			client, exists := clients[platform.Name]
			if !exists {
				logger.LogError(fmt.Sprintf("Не найден клиент для %s", platform.Name))
				return
			}

			if err = client.Upload(videoPath, content[0].Title, content[0].Excerpt, []string{"tags"}); err != nil {
				logger.LogError(fmt.Sprintf("Ошибка публикации на платформе %s: %v", platform.Name, err))
			}
		}(platform)
	}

	wg.Wait()

	logger.LogInfo(fmt.Sprint("Видео успешно обработано", "email", user.Email, "path", videoPath))

	return nil
}

// GenerateVideo - основная логика генерации видео
func GenerateVideo(text, audioPath string) (string, error) {
	logger.LogInfo(fmt.Sprint("Генерация видео", "text", text))

	textVideoPath, err := generateTextVideo(text)
	if err != nil {
		return "", errors.Wrap(err, "ошибка генерации текста в видео")
	}

	videoPath, err := combineAudioWithVideo(textVideoPath, audioPath)
	if err != nil {
		return "", errors.Wrap(err, "ошибка наложения аудио")
	}

	if err = removeFile(textVideoPath); err != nil {
		logger.LogError(fmt.Sprint("Не удалось удалить временное видео", "file", textVideoPath))
	}

	return videoPath, nil
}

func generateTextVideo(text string) (string, error) {
	videoPath := fmt.Sprintf("/tmp/%d_text_video.mp4", time.Now().Unix())

	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-t", "5", "-i", "color=c=black:s=1280x720:r=30",
		"-vf", fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=48:x=(w-text_w)/2:y=(h-text_h)/2", text),
		"-an", videoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка ffmpeg: %v, output: %s", err, string(output))
	}

	return videoPath, nil
}

func GenerateAudioForText(text string) (string, error) {
	audioPath, err := speech.GenerateAudio(text)
	if err != nil {
		return "", err
	}
	return audioPath, nil
}

func combineAudioWithVideo(videoPath, audioPath string) (string, error) {
	finalVideoPath := fmt.Sprintf("/tmp/%d_final_video.mp4", time.Now().Unix())

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-c:v", "libx264", "-c:a", "aac",
		"-strict", "experimental", "-shortest", finalVideoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка объединения видео и аудио: %v, output: %s", err, string(output))
	}

	return finalVideoPath, nil
}

func removeFile(path string) error {
	if err := os.Remove(path); err != nil {
		logger.LogError(fmt.Sprint("Ошибка удаления файла", "file", path, "error", err))
		return err
	}
	return nil
}
