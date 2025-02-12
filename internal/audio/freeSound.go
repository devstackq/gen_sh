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

type SearchResponse struct {
	Count    int     `json:"count"`
	Previous *string `json:"previous"` // Nullable, use *string for null values
	Next     *string `json:"next"`     // Nullable, use *string for null values
	Results  []Sound `json:"results"`
}

type Sound struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	License  string   `json:"license"`
	Username string   `json:"username"`
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

	var result SearchResponse

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON Freesound: %v", err)
	}

	audioResults := make([]AudioResult, 0, len(result.Results))
	for _, r := range result.Results {
		audioResults = append(audioResults, AudioResult{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	if len(audioResults) == 0 {
		return nil, fmt.Errorf("sounds result is empty")
	}

	//DRY ? or move func?
	soundURL := fmt.Sprintf("https://freesound.org/apiv2/sounds/%d/", audioResults[0].ID)

	req, err := http.NewRequest(http.MethodGet, soundURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", "Token "+f.ApiKey)

	client := &http.Client{}
	soundResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer soundResp.Body.Close()

	defer soundResp.Body.Close()

	soundBody, err := io.ReadAll(soundResp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа Freesound: %v", err)
	}

	fmt.Println(string(soundBody), "1 sound")

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
