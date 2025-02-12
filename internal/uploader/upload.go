package uploader

import (
	"fmt"

	"github.com/devstackq/gen_sh/internal/config"
)

type PlatformClient interface {
	Upload(videoPath, title, description string, tags []string) error
}

func New(platformConfig config.Platform) (PlatformClient, error) {
	switch platformConfig.Name {
	case "youtube":
		return NewYtClient(platformConfig)
	case "tiktok":
		//return NewTikTokClient(platformConfig)
	}
	return nil, fmt.Errorf("unknown platform %s", platformConfig.Name)
}
