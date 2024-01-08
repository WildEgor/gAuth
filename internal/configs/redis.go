package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type RedisConfig struct {
	URI string `env:"REDIS_URI,required"`
}

func NewRedisConfig(c *Configurator) *RedisConfig {
	cfg := RedisConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
