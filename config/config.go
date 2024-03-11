package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DbHost     string `env:"DB_HOSTNAME" env-default:"localhost"`
	DbPort     int    `env:"DB_PORT" env-default:"5432"`
	DbUser     string `env:"DB_USERNAME" env-default:"postgres"`
	DbPassword string `env:"DB_PASSWORD" env-default:"root"`
	DbName     string `env:"DB_NAME" env-default:"assignment2"`
}

func LoadEnv() EnvConfig {
	cfg := EnvConfig{}

	err := godotenv.Load()
	if err != nil {
		panic("cannot load env file")
	}

	err = env.Parse(&cfg)
	if err != nil {
		panic("cannot load env file")
	}
	return cfg
}
