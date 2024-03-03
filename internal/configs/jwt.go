package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
	"time"
)

type JWTConfig struct {
	Secret string `env:"JWT_SECRET,required"`
	at     string `env:"JWT_AT_TTL" envDefault:"10m"`
	rt     string `env:"JWT_RT_TTL" envDefault:"24h"`

	ATDuration time.Duration `env:"-"`
	RTDuration time.Duration `env:"-"`
}

func NewJWTConfig(c *Configurator) *JWTConfig {
	cfg := JWTConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	cfg.ATDuration, _ = time.ParseDuration(cfg.at)
	cfg.RTDuration, _ = time.ParseDuration(cfg.rt)

	return &cfg
}
