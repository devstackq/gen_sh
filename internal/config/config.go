package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Platform struct {
	Name        string `yaml:"name"`
	Credentials string `yaml:"credentials"`
	APIKey      string `yaml:"api_key"`
	UploadPath  string `yaml:"upload_path"`
}

type User struct {
	Email     string     `yaml:"email"`
	Theme     string     `yaml:"theme"`
	Sources   []string   `yaml:"sources"`
	Platforms []Platform `yaml:"platforms"`
}

type Config struct {
	Users []User `yaml:"users"`
}

// LoadConfig загружает конфигурацию из YAML файла
func LoadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть конфигурационный файл: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфигурацию: %v", err)
	}

	return &config, nil
}
