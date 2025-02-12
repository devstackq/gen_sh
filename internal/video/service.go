package video

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/content"
	"github.com/devstackq/gen_sh/internal/logger"
	"github.com/devstackq/gen_sh/internal/speech"
	"github.com/devstackq/gen_sh/internal/stock"
	"github.com/devstackq/gen_sh/internal/uploader"
	"github.com/pkg/errors"
)

type Video struct {
	uploader uploader.PlatformClient
}

func Publish(user config.User, item content.Content) error {

	logger.LogInfo(fmt.Sprint("Начата обработка пользователя", "email", user.Email, "theme", user.Theme))

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

			if err := client.Upload(item.Path, item.Title, item.Description, item.Tags); err != nil { //todo gen - tags
				logger.LogError(fmt.Sprintf("Ошибка публикации на платформе %s: %v", platform.Name, err))
			}
		}(platform)
	}

	wg.Wait()

	logger.LogInfo(fmt.Sprint("Видео успешно обработано", "email", user.Email, "path", item.Path))

	return nil
}

func GenerateVideo(user config.User, content []content.Content) (string, error) {

	var (
		mediaType = "video" // photo/video - getFromConfig?
		perPage   = 1       // getFromConfig?
		text      = content[0].Excerpt
	)

	stock := stock.New("pexels")

	medias, err := stock.SearchMedia(user.Theme, mediaType, perPage)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска медиафайлов %v", err)
	}

	if len(medias) == 0 {
		return "", fmt.Errorf("не найдено подходящих медиафайлов")
	}

	videoURL := medias[0].Source
	videoPath, err := downloadVideo(videoURL)
	if err != nil {
		return "", errors.Wrap(err, "ошибка загрузки видео")
	}
	defer os.Remove(videoPath)

	audioPath, err := speech.GenerateAudio(text)
	if err != nil {
		return "", err
	}

	//title, description, tags := GenerateMetadata(user.Theme)
	logger.LogInfo(fmt.Sprint("Генерация видео", "text", text))

	finalVideoPath, err := combineAudioWithVideo(videoPath, audioPath, "path/to/watermark.png") //todo set logo image
	if err != nil {
		return "", errors.Wrap(err, "ошибка наложения аудио")
	}

	if err = removeFile(videoPath); err != nil {
		logger.LogError(fmt.Sprint("Не удалось удалить временное видео", "file", videoPath))
	}

	return finalVideoPath, nil
}

func downloadVideo(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("не удалось загрузить видео, статус: %s", resp.Status)
	}

	tempDir := os.TempDir()
	fileName := fmt.Sprintf("%d_video.mp4", time.Now().Unix())
	filePath := filepath.Join(tempDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения видео: %v", err)
	}

	return filePath, nil
}

func generateTextVideo(text string) (string, error) {
	videoPath := fmt.Sprintf("/tmp/%d_text_video.mp4", time.Now().Unix())

	cmd := exec.Command("ffmpeg", "-f", "lavfi", "-t", "30", "-i", "color=c=black:s=1280x720:r=30",
		"-vf", fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=48:x=(w-text_w)/2:y=(h-text_h)/2", text),
		"-an", videoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка ffmpeg: %v, output: %s", err, string(output))
	}

	return videoPath, nil
}

func combineAudioWithVideo(videoPath, audioPath, watermarkPath string) (string, error) {
	finalVideoPath := fmt.Sprintf("/tmp/%d_final_video.mp4", time.Now().Unix())

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-i", audioPath, "-i", watermarkPath,
		"-filter_complex", "[0:v][2:v]overlay=W-w-10:H-h-10:format=auto[v]", "-map", "[v]", "-map", "1:a",
		"-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-shortest", finalVideoPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка объединения видео, аудио и водяного знака: %v, output: %s", err, string(output))
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
