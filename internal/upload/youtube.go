package upload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/devstackq/gen_sh/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/youtube/v3"
)

// YouTubeUploader — структура для загрузки видео на YouTube
type YouTubeUploader struct {
	client *youtube.Service
}

// NewYouTubeUploader — создаёт нового клиента для YouTube API
func NewYouTubeUploader(credentialsFile string) (*YouTubeUploader, error) {
	// Загружаем credentials из файла
	file, err := os.Open(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии credentials: %v", err)
	}
	defer file.Close()

	config, err := google.ConfigFromJSON(file, youtube.YoutubeUploadScope)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении конфигурации OAuth2: %v", err)
	}

	// Получаем токен OAuth2
	client := getClient(config)

	// Создаем клиент YouTube API
	youtubeService, err := youtube.New(client)
	if err != nil {
		return nil, fmt.Errorf("ошибка при инициализации YouTube API: %v", err)
	}

	return &YouTubeUploader{
		client: youtubeService,
	}, nil
}

// getClient — функция для аутентификации с OAuth2
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile("token.json")
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken("token.json", tok)
	}
	return config.Client(oauth2.NoContext, tok)
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

// getTokenFromWeb — получает токен через веб-интерфейс
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по следующему URL и введите код: \n%v\n", authURL)

	var code string
	fmt.Print("Введите код: ")
	fmt.Scan(&code)

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("Ошибка при получении токена: %v", err)
	}
	return tok
}

// saveToken — сохраняет токен в файл
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Сохраняем токен в %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		logger.Log.Fatalf("Ошибка при создании файла токена: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Upload — метод загрузки видео на YouTube
func (u *YouTubeUploader) Upload(videoPath string) error {
	// Открываем видеофайл
	file, err := os.Open(videoPath)
	if err != nil {
		logger.Log.Error("Не удалось открыть видеофайл", zap.String("path", videoPath))
		return fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Создаем запрос на загрузку
	call := u.client.Videos.Insert([]string{"snippet", "status"}, &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       "Название видео",
			Description: "Описание видео",
			Tags:        []string{"тег1", "тег2"},
		},
		Status: &youtube.VideoStatus{
			PrivacyStatus: "public", // публичное видео
		},
	}, file)

	// Отправляем запрос
	_, err = call.Do()
	if err != nil {
		logger.Log.Error("Ошибка при загрузке видео на YouTube", zap.String("path", videoPath), zap.Error(err))
		return fmt.Errorf("не удалось загрузить видео: %v", err)
	}

	logger.Log.Info("✅ Видео успешно загружено на YouTube", zap.String("path", videoPath))
	return nil
}
