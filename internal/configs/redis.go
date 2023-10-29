package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type RedisConfig struct {
	URI string `env:"REDIS_URI"`
}

func NewRedisConfig(c *Configurator) *RedisConfig {
	cfg := RedisConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[RedisConfig] %+v\n", err)
	}

	return &cfg
}
