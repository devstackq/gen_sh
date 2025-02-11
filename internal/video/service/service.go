package service

import (
	"github.com/devstackq/gen_sh/internal/video/model"
	"github.com/devstackq/gen_sh/internal/video/repository"
)

type VideoService interface {
	CreateVideo(title, description, url string) (*model.Video, error)
	GetVideoByID(id int64) (*model.Video, error)
	GetAllVideos() ([]*model.Video, error)
}

type videoService struct {
	repo repository.VideoRepository
}

func NewVideoService(repo repository.VideoRepository) VideoService {
	return &videoService{repo: repo}
}

func (s *videoService) CreateVideo(title, description, url string) (*model.Video, error) {
	video := &model.Video{
		Title:       title,
		Description: description,
		URL:         url,
	}
	return s.repo.Create(video)
}

func (s *videoService) GetVideoByID(id int64) (*model.Video, error) {
	return s.repo.GetByID(id)
}

func (s *videoService) GetAllVideos() ([]*model.Video, error) {
	return s.repo.GetAll()
}
