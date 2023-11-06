package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type KeycloakConfig struct {
	ClientID     string `env:"KEYCLOAK_CLIENT_ID"`
	ClientSecret string `env:"KEYCLOAK_CLIENT_SECRET"`
	RSAPublicKey string `env:"KEYCLOAK_CLIENT_RSA_PUBLIC_KEY"`
	Realm        string `env:"KEYCLOAK_REALM"`
	API          string `env:"KEYCLOAK_URL"`
}

func NewKeycloakConfig() *KeycloakConfig {
	cfg := KeycloakConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[KeycloakConfig] %+v\n", err)
	}

	return &cfg
}
