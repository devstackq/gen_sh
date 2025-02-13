package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/devstackq/gen_sh/internal/config"
)

const (
	apiBaseURL = "https://api.pexels.com/v1"
)

type pexels struct{}

func (p *pexels) SearchMedia(user config.User, mediaType string, perPage int, duration float64) ([]MediaItem, error) {

	searchURL := fmt.Sprintf("%s/search?query=%s&per_page=%d", apiBaseURL, url.QueryEscape(user.Theme), perPage)

	if mediaType == "video" {
		searchURL = fmt.Sprintf(
			"%s/videos/search?query=%s&per_page=%d&min_duration=%f&max_duration=%f",
			apiBaseURL, url.QueryEscape(user.Theme), perPage, duration-5, duration+5,
		)
	}

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", user.Stock.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось выполнить запрос: статус %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mediaItems []MediaItem
	if mediaType == "video" {
		var result struct {
			Videos []struct {
				ID         int    `json:"id"`
				URL        string `json:"url"`
				VideoFiles []struct {
					Link string `json:"link"`
				} `json:"video_files"`
			} `json:"videos"`
		}
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		for _, video := range result.Videos {
			mediaItems = append(mediaItems, MediaItem{
				ID:     video.ID,
				Type:   "video",
				URL:    video.URL,
				Source: video.VideoFiles[0].Link,
			})
		}
	} else {
		var result struct {
			Photos []struct {
				ID  int    `json:"id"`
				URL string `json:"url"`
				Src struct {
					Original string `json:"original"`
				} `json:"src"`
			} `json:"photos"`
		}
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		for _, photo := range result.Photos {
			mediaItems = append(mediaItems, MediaItem{
				ID:     photo.ID,
				Type:   "photo",
				URL:    photo.URL,
				Source: photo.Src.Original,
			})
		}
	}

	fmt.Println(mediaItems, "media")

	return mediaItems, nil
}
