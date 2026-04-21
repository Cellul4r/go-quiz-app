package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	AppDebug bool   `mapstructure:"APP_DEBUG"`
	Port     string `mapstructure:"PORT"`
	DBUrl    string `mapstructure:"DB_URL"`
}

var envs = []string{
	"APP_ENV",
	"APP_DEBUG",
	"PORT",
	"DB_URL",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("/")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
