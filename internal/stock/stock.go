package stock

type Stock interface {
	SearchMedia(query string, mediaType string, perPage int) ([]MediaItem, error)
}

type MediaItem struct {
	ID     int    `json:"id"`
	Type   string `json:"type"` // "photo" или "video"
	URL    string `json:"url"`
	Source string `json:"src"` // Ссылка на файл
}

func New(source string) *pexels {
	if source == "pexels" {
		return &pexels{}
	}
	return nil
}
