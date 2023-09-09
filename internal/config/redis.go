package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type RedisConfig struct {
	URI string `env:"REDIS_URI"`
}

func NewRedisConfig() *RedisConfig {
	cfg := RedisConfig{}

	if err := godotenv.Load(".env"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[RedisConfig] %+v\n", err)
		}
	}

	return &cfg
}
