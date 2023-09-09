package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type KeycloakConfig struct {
	API string `env:"KEYCLOAK_URL"`
}

func NewKeycloakConfig() *KeycloakConfig {
	cfg := KeycloakConfig{}

	if err := godotenv.Load(".env"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[KeycloakConfig] %+v\n", err)
		}
	}

	return &cfg
}
