package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
	AT     int64  `env:"JWT_AT_TTL" envDefault:"3600"`
	RT     int64  `env:"JWT_RT_TTL" envDefault:"86400"`

	ATDuration time.Duration `env:"-"`
	RTDuration time.Duration `env:"-"`
}

func NewJWTConfig(c *Configurator) *JWTConfig {
	cfg := JWTConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	cfg.ATDuration = time.Duration(cfg.AT) * time.Second
	cfg.RTDuration = time.Duration(cfg.RT) * time.Second

	return &cfg
}
