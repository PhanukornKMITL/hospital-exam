package config

import "os"

type Config struct {
	AppPort string
	AppEnv  string
}

func Load() *Config {
	return &Config{
		AppPort: os.Getenv("APP_PORT"),
		AppEnv:  os.Getenv("APP_ENV"),
	}
}
