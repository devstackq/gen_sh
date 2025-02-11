package repo

import (
	"database/sql"
	"github.com/devstackq/gen_sh/internal/video/model"
	"log"
)

type VideoRepository interface {
	Create(video *model.Video) (*model.Video, error)
	GetByID(id int64) (*model.Video, error)
	GetAll() ([]*model.Video, error)
}

type videoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) Create(video *model.Video) (*model.Video, error) {
	query := `INSERT INTO videos (title, description, url) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, video.Title, video.Description, video.URL).Scan(&video.ID, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		log.Println("Error creating video:", err)
		return nil, err
	}
	return video, nil
}

func (r *videoRepository) GetByID(id int64) (*model.Video, error) {
	query := `SELECT id, title, description, url, created_at, updated_at FROM videos WHERE id = $1`
	var video model.Video
	err := r.db.QueryRow(query, id).Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CreatedAt, &video.UpdatedAt)
	if err != nil {
		log.Println("Error fetching video by ID:", err)
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) GetAll() ([]*model.Video, error) {
	query := `SELECT id, title, description, url, created_at, updated_at FROM videos`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error fetching videos:", err)
		return nil, err
	}
	defer rows.Close()

	var videos []*model.Video
	for rows.Next() {
		var video model.Video
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.URL, &video.CreatedAt, &video.UpdatedAt); err != nil {
			log.Println("Error scanning video:", err)
			return nil, err
		}
		videos = append(videos, &video)
	}
	return videos, nil
}
