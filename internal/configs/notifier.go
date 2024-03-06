package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type NotifierConfig struct {
	DSN      string `env:"AMQP_DSN,required"`
	Exchange string `env:"NOTIFIER_EXCHANGE,required"`
}

func NewNotifierConfig(c *Configurator) *NotifierConfig {
	cfg := NotifierConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
