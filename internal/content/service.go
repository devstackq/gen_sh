package content

import (
	"fmt"
	"strings"
)

// NewContentFetcher – фабричная функция для создания нужного fetcher-а по источнику.
func NewContentFetcher(source string) (Fetcher, error) {
	switch strings.ToLower(source) {
	case "reddit":
		return &RedditFetcher{}, nil
	case "wikipedia":
		return &WikipediaFetcher{}, nil
	case "twitter":
		return &TwitterFetcher{}, nil
	default:
		return nil, fmt.Errorf("неизвестный источник: %s", source)
	}
}

func FetchContent(theme string, sources []string) ([]Content, error) {
	var allItems []Content

	for _, src := range sources {
		fetcher, err := NewContentFetcher(src)
		if err != nil {
			// Пропускаем неизвестные источники, можно логировать ошибку
			continue
		}
		items, err := fetcher.Fetch(theme)
		if err != nil {
			// Логируем ошибку, но продолжаем
			continue
		}
		allItems = append(allItems, items...)
	}

	return allItems, nil
}
