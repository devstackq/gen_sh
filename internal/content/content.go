package content

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Content struct {
	Source      string // Например, "Reddit", "Wikipedia", "Twitter"
	Title       string
	Description string //for upload
	URL         string
	Excerpt     string   // Краткий отрывок или выдержка
	Text        string   // Полное описание или текст статьи/поста
	Tags        []string // Теги, сгенерированные на основе заголовка или анализа текста

	Path string
}

type Fetcher interface {
	Fetch(theme string) ([]Content, error)
}

type RedditFetcher struct{}

func (rf *RedditFetcher) Fetch(theme string) ([]Content, error) {
	subreddit := url.QueryEscape(strings.ToLower(theme))
	apiURL := fmt.Sprintf("https://www.reddit.com/r/%s/top/.json?limit=5&t=day", subreddit)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ContentFetcherBot/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var redditResp RedditResponse
	if err = json.Unmarshal(body, &redditResp); err != nil {
		return nil, err
	}

	var items []Content
	for _, child := range redditResp.Data.Children {
		// Если selftext пустой, используем заголовок в качестве полного текста.
		fullText := child.Data.Selftext
		if fullText == "" {
			fullText = child.Data.Title
		}
		// Генерация тегов на основе заголовка.
		tags := generateTags(child.Data.Title)
		item := Content{
			Source:  "Reddit",
			Title:   child.Data.Title,
			URL:     child.Data.URL,
			Excerpt: child.Data.Selftext,
			Text:    fullText,
			Tags:    tags,
		}
		items = append(items, item)
	}
	fmt.Printf("Found %d reddit items.\n", len(items))

	return items, nil
}

// RedditResponse описывает структуру ответа Reddit API.
type RedditResponse struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Kind string `json:"kind"`
			Data struct {
				Title    string `json:"title"`
				URL      string `json:"url"`
				Selftext string `json:"selftext"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type WikipediaFetcher struct{}

func (wf *WikipediaFetcher) Fetch(theme string) ([]Content, error) {
	escapedTheme := url.QueryEscape(theme)
	apiURL := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/page/summary/%s", escapedTheme)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from Wikipedia API: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var wpResp WikipediaResponse
	if err = json.Unmarshal(body, &wpResp); err != nil {
		return nil, err
	}

	// Используем Extract как полный текст и краткий отрывок.
	tags := []string{"wikipedia", strings.ToLower(wpResp.Title)}
	item := Content{
		Source:  "Wikipedia",
		Title:   wpResp.Title,
		URL:     wpResp.ContentUrls.Desktop.Page,
		Excerpt: wpResp.Extract,
		Text:    wpResp.Extract,
		Tags:    tags,
	}
	return []Content{item}, nil
}

type WikipediaResponse struct {
	Title       string `json:"title"`
	Extract     string `json:"extract"`
	ContentUrls struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

type TwitterFetcher struct{}

func (tf *TwitterFetcher) Fetch(theme string) ([]Content, error) {
	tags := []string{"twitter", strings.ToLower(theme)}
	item := Content{
		Source:  "Twitter",
		Title:   fmt.Sprintf("Пример твита по теме %s", theme),
		URL:     "https://twitter.com/example",
		Excerpt: "Это пример содержимого твита.",
		Text:    "Это пример содержимого твита.",
		Tags:    tags,
	}
	return []Content{item}, nil
}

func generateTags(title string) []string {
	words := strings.Fields(title)
	var tags []string
	for _, word := range words {
		clean := strings.ToLower(strings.Trim(word, ".,;:\"'!?"))
		if len(clean) > 3 {
			tags = append(tags, clean)
		}
	}
	return tags
}
