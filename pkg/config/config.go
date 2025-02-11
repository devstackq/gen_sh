package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

var Cfg Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфигурации: %v", err)
	}

	Cfg.DatabaseURL = viper.GetString("database.url")
}
