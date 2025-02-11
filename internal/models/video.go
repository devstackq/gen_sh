package models

import "time"

// Video — структура для хранения информации о видео
type Video struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
	Platform  string    `json:"platform"`
	Uploaded  bool      `json:"uploaded"`
	CreatedAt time.Time `json:"created_at"`
}
