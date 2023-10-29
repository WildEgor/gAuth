package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type KeycloakConfig struct {
	API string `env:"KEYCLOAK_URL"`
}

func NewKeycloakConfig(c *Configurator) *KeycloakConfig {
	cfg := KeycloakConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[KeycloakConfig] %+v\n", err)
	}

	return &cfg
}