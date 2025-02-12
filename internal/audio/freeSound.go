package audio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const freesoundAPI = "https://freesound.org/apiv2/search/text/"

type FreeSoundClient struct {
	ApiKey string
}

type AudioResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"previews"`
}

type SoundResult struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	License     string            `json:"license"`
	Previews    map[string]string `json:"previews"`
	Duration    float64           `json:"duration"`
	DownloadURL string            `json:"download"`
	Score       float64           `json:"score,omitempty"`
}

type SoundListResponse struct {
	Results []SoundResult `json:"results"`
}

func NewFreeSoundClient(apiKey string) *FreeSoundClient {
	return &FreeSoundClient{ApiKey: apiKey}
}

func (f *FreeSoundClient) Search(query string, limit int) ([]AudioResult, error) {
	url := fmt.Sprintf("%s?q=%s&token=%s&page_size=%d", freesoundAPI, query, f.ApiKey, limit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к API Freesound: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа Freesound: %v", err)
	}

	var result SoundListResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON Freesound: %v", err)
	}

	audioResults := make([]AudioResult, 0, len(result.Results))
	for _, r := range result.Results {
		if previewURL, ok := r.Previews["preview-hq-mp3"]; ok {
			audioResults = append(audioResults, AudioResult{
				ID:   r.ID,
				Name: r.Name,
				URL:  previewURL,
			})
		}
	}

	return audioResults, nil
}

func DownloadAudio(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("ошибка загрузки аудиофайла: %v", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи аудиофайла: %v", err)
	}

	return nil
}
