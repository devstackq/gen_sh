package uploader

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/devstackq/gen_sh/internal/config"
	"github.com/devstackq/gen_sh/internal/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type ytService struct {
	client *youtube.Service
}

// NewYtClient создает новый объект загрузчика YouTube
func NewYtClient(ytConfig config.Platform) (ytService, error) {
	logger.LogInfo("Инициализация YouTube uploader...")
	ctx := context.Background()

	data, err := os.ReadFile(ytConfig.Credentials)
	if err != nil {
		return ytService{}, fmt.Errorf("ошибка при чтении credentials: %v", err)
	}

	conf, err := google.ConfigFromJSON(data, youtube.YoutubeUploadScope)
	if err != nil {
		return ytService{}, fmt.Errorf("ошибка при получении конфигурации OAuth2: %v", err)
	}

	// Получаем OAuth клиент
	client, err := getClient(ctx, conf)
	if err != nil {
		return ytService{}, err
	}
	// Создаем клиент YouTube API
	youtubeService, err := youtube.New(client)
	if err != nil {
		return ytService{}, fmt.Errorf("ошибка при инициализации YouTube API: %v", err)
	}

	return ytService{client: youtubeService}, nil
}

// getClient — функция для аутентификации с OAuth2
func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, tok), nil //oauth2.NoContext used
}

// tokenFromFile — загружает токен OAuth из файла
func tokenFromFile(file string) (*oauth2.Token, error) {
	tokFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer tokFile.Close()

	var token oauth2.Token
	err = json.NewDecoder(tokFile).Decode(&token)
	return &token, err
}

func (u *ytService) Upload(videoPath, title, description string, tags []string) error {
	logger.LogInfo(fmt.Sprintf("Загрузка видео на YouTube: %s", videoPath))

	// Открываем видеофайл
	file, err := os.Open(videoPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Создаем запрос на загрузку
	call := u.client.Videos.Insert(
		[]string{"snippet", "status"},
		&youtube.Video{
			Snippet: &youtube.VideoSnippet{
				Title:       title,
				Description: description,
				Tags:        tags,
			},
			Status: &youtube.VideoStatus{
				PrivacyStatus: "public",
			},
		})

	call = call.Media(file)

	_, err = call.Do()
	if err != nil {
		return fmt.Errorf("не удалось загрузить видео: %v", err)
	}

	logger.LogInfo(fmt.Sprintf("✅ Видео успешно загружено на YouTube: %s", videoPath))
	return nil
}
