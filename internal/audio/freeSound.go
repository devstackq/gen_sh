package audio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const baseURL = "https://freesound.org"

type FreeSoundClient struct {
	ApiKey string
}

type AudioResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
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

func (f *FreeSoundClient) Search(query string, limit int, duration float64) (AudioResult, error) {
	var result AudioResult

	filter := url.QueryEscape("duration:[8.6 TO 18.6]")
	searchURL := fmt.Sprintf("%s?q=%s&filter=%s&token=%s&page_size=%d",
		fmt.Sprint(baseURL+"/apiv2/search/text/"), url.QueryEscape(query), filter, f.ApiKey, limit)

	searchResp, err := http.Get(searchURL)
	if err != nil {
		return result, fmt.Errorf("ошибка запроса к API Freesound: %v", err)
	}
	defer searchResp.Body.Close()

	searchBody, err := io.ReadAll(searchResp.Body)
	if err != nil {
		return result, fmt.Errorf("ошибка чтения ответа Freesound: %v", err)
	}

	if searchResp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("statusCode is not OK")
	}

	var searchResult SearchResponse
	if err = json.Unmarshal(searchBody, &searchResult); err != nil {
		return result, fmt.Errorf("ошибка парсинга JSON Freesound: %v", err)
	}

	audioResults := make([]AudioResult, 0, len(searchResult.Results))
	for _, r := range searchResult.Results {
		audioResults = append(audioResults, AudioResult{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	if len(audioResults) == 0 {
		return result, fmt.Errorf("sounds result is empty")
	}

	//DRY ? or move func?
	soundURL := fmt.Sprintf("https://freesound.org/apiv2/sounds/%d/", audioResults[0].ID)

	req, err := http.NewRequest(http.MethodGet, soundURL, nil)
	if err != nil {
		return result, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Authorization", "Token "+f.ApiKey)

	client := &http.Client{}
	soundResp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer soundResp.Body.Close()

	soundBody, err := io.ReadAll(soundResp.Body)
	if err != nil {
		return result, fmt.Errorf("ошибка чтения ответа Freesound: %v", err)
	}

	if err = json.Unmarshal(soundBody, &result); err != nil {
		return result, fmt.Errorf("ошибка парсинга JSON Freesound: %v", err)
	}
	fmt.Println(result, "1 sound")

	return result, nil
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
