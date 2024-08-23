package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort string
}

// nolint: gochecknoglobals
var configInstance *Config

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("can't load config from `.env`. environment variables will be used. err: %v", err)
	}

	configInstance = &Config{
		AppPort: viper.GetString("APP_PORT"),
	}
}

func Instance() *Config {
	return configInstance
}
