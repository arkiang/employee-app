package configs

import (
	"time"

	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin   string `mapstructure:"CLIENT_ORIGIN"`

	JwtSecret      string        `mapstructure:"JWT_SECRET_KEY"`
	JWTExpiration  time.Duration `mapstructure:"JWT_EXPIRATION_MINUTES"`
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&AppConfig)
	return
}