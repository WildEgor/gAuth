package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Port    string `env:"APP_PORT" envDefault:"8888"`
	RPCPort string `env:"GRPC_PORT" envDefault:"8887"`
	Mode    string `env:"APP_MODE,required"`
	GoEnv   string `env:"GO_ENV" envDefault:"local"`
	Version string `env:"VERSION" envDefault:"local"`
}

func NewAppConfig(c *Configurator) *AppConfig {
	cfg := AppConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}

func (ac AppConfig) IsProduction() bool {
	return ac.Mode != "develop"
}
