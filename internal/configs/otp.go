package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type OTPConfig struct {
	Issuer string `env:"OTP_ISSUER,required"`
	Secret string `env:"OTP_SECRET,required"`
}

func NewOTPConfig(c *Configurator) *OTPConfig {
	cfg := OTPConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
