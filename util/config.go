package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBName     string `mapstructure:"DB_NAME"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`

	SessionName       string `mapstructure:"SESSION_NAME"`
	ExpiryInSecond    int    `mapstructure:"EXPIRY_IN_SECOND"`
	IsSessionSecure   bool   `mapstructure:"IS_SESSION_SECURE"`
	IsSessionHttpOnly bool   `mapstructure:"IS_SESSION_HTTP_ONLY"`

	BaseUrl  string `mapstructure:"BASE_URL"`
	BasePath string `mapstructure:"BASE_PATH"`
	PORT     int    `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
