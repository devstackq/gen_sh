package upload

import (
	"fmt"
	"github.com/devstackq/gen_sh/pkg/logger"
)

// VideoUploader — интерфейс для загрузки видео на различные платформы
type VideoUploader interface {
	Upload(videoPath string) error
}
