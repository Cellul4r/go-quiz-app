package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv            string `mapstructure:"APP_ENV"`
	AppDebug          bool   `mapstructure:"APP_DEBUG"`
	Port              string `mapstructure:"PORT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	SupabaseURL       string `mapstructure:"SUPABASE_URL"`
	SupabaseJWTSecret string `mapstructure:"SUPABASE_JWT_SECRET"`
}

var envs = []string{
	"APP_ENV",
	"APP_DEBUG",
	"PORT",
	"DB_HOST",
	"DB_PORT",
	"DB_NAME",
	"DB_USER",
	"DB_PASSWORD",
	"SUPABASE_URL",
	"SUPABASE_JWT_SECRET",
}

func LoadConfig() (Config, error) {
	var config Config

	for _, env := range envs {
		viper.BindEnv(env)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
