package uploader

import (
	"fmt"
	"github.com/devstackq/gen_sh/internal/config"
)

// PlatformClient - интерфейс для всех платформ
type PlatformClient interface {
	Upload(videoPath, title, description string, tags []string) error
}

func New(platformConfig config.Platform) (PlatformClient, error) {
	switch platformConfig.Name {
	case "youtube":
		return NewYtClient(platformConfig), nil
	case "tiktok":
		//return NewTikTokClient(platformConfig), nil
	}
	return nil, fmt.Errorf("unknown platform %s", platformConfig.Name)
}
