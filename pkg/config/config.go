package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
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
		AppPort:    viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBPort:     viper.GetString("DB_PORT"),
	}
}

func Instance() *Config {
	return configInstance
}
