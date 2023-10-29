package keycloak_adapter

import (
	"github.com/Nerzal/gocloak/v3"
	"github.com/WildEgor/gAuth/internal/configs"
	log "github.com/sirupsen/logrus"
)

type KeycloakAdapter struct {
	Client         gocloak.GoCloak
	keycloakConfig *configs.KeycloakConfig
}

// NewKeycloakAdapter Create new Keycloak Adapter
func NewKeycloakAdapter(
	keycloakConfig *configs.KeycloakConfig,
) *KeycloakAdapter {
	adapter := &KeycloakAdapter{
		nil,
		keycloakConfig,
	}

	adapter.newClient()

	return adapter
}

func (ka *KeycloakAdapter) newClient() {
	ka.Client = gocloak.NewClient(ka.keycloakConfig.API)
	log.Info("Success init Keycloak")
}

// Login Auth in Keycloak
func (ka *KeycloakAdapter) Login(user string, pass string) (*JWT, error) {
	token, err := ka.Client.Login("admin-cli", "", "master", user, pass)
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
