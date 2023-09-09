package keycloak_adapter

import (
	"github.com/Nerzal/gocloak/v3"
	"github.com/WildEgor/gAuth/internal/config"
	log "github.com/sirupsen/logrus"
)

var (
	client gocloak.GoCloak
)

type KeycloakAdapter struct {
	keycloakConfig *config.KeycloakConfig
}

func NewKeycloakAdapter(
	keycloakConfig *config.KeycloakConfig,
) *KeycloakAdapter {
	adapter := &KeycloakAdapter{
		keycloakConfig,
	}

	adapter.newClient()

	return adapter
}

func (ka *KeycloakAdapter) newClient() {
	keycloakClient := gocloak.NewClient(ka.keycloakConfig.API)
	client = keycloakClient
	log.Info("Success init Keycloak")
}

func (ka *KeycloakAdapter) Login(user string, pass string) (*JWT, error) {
	token, err := client.Login("admin-cli", "", "master", user, pass)
	if err != nil {
		log.Error("Keycloak login failed")
		return nil, err
	}

	var jwt *JWT
	jwt.AccessToken = token.AccessToken
	jwt.IDToken = token.IDToken
	jwt.ExpiresIn = token.ExpiresIn
	jwt.RefreshExpiresIn = token.RefreshExpiresIn
	jwt.RefreshToken = token.RefreshToken
	jwt.TokenType = token.TokenType
	jwt.NotBeforePolicy = token.NotBeforePolicy
	jwt.SessionState = token.SessionState
	jwt.Scope = token.Scope

	return jwt, nil
}
