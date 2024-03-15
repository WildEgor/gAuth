package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type OTPConfig struct {
	Issuer string `env:"OTP_ISSUER" envDefault:"gAuth"`
	Length uint8  `env:"OTP_LENGTH" envDefault:"6"`
}

func NewOTPConfig(c *Configurator) *OTPConfig {
	cfg := OTPConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
