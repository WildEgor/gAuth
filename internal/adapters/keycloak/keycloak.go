package keycloak_adapter

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type KeycloakAdapter struct {
	Client         *gocloak.GoCloak
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

func (ka *KeycloakAdapter) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	token, err := ka.Client.LoginClient(ctx, ka.keycloakConfig.ClientID, ka.keycloakConfig.ClientSecret, ka.keycloakConfig.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login the rest client")
	}
	return token, nil
}

func (ka *KeycloakAdapter) CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error) {
	token, err := ka.loginClient(ctx)
	if err != nil {
		return nil, err
	}

	userId, err := ka.Client.CreateUser(ctx, token.AccessToken, ka.keycloakConfig.Realm, user)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create the user")
	}

	err = ka.Client.SetPassword(ctx, token.AccessToken, userId, ka.keycloakConfig.Realm, password, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}

	var roleNameLowerCase = strings.ToLower(role)
	roleKeycloak, err := ka.Client.GetRealmRole(ctx, token.AccessToken, ka.keycloakConfig.Realm, roleNameLowerCase)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", roleNameLowerCase))
	}
	err = ka.Client.AddRealmRoleToUser(ctx, token.AccessToken, ka.keycloakConfig.Realm, userId, []gocloak.Role{
		*roleKeycloak,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add a realm role to user")
	}

	userKeycloak, err := ka.Client.GetUserByID(ctx, token.AccessToken, ka.keycloakConfig.Realm, userId)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get recently created user")
	}

	return userKeycloak, nil
}

// Login Auth in Keycloak
func (ka *KeycloakAdapter) Login(user string, pass string) (*JWT, error) {
	token, err := ka.Client.Login(context.Background(), ka.keycloakConfig.ClientID, ka.keycloakConfig.ClientSecret, ka.keycloakConfig.Realm, user, pass)
	if err != nil {
		log.Error("Keycloak login failed")
		return nil, err
	}

	var jwt *JWT = &JWT{
		AccessToken:      token.AccessToken,
		IDToken:          token.IDToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		RefreshToken:     token.RefreshToken,
		TokenType:        token.TokenType,
		NotBeforePolicy:  token.NotBeforePolicy,
		SessionState:     token.SessionState,
		Scope:            token.Scope,
	}

	return jwt, nil
}
