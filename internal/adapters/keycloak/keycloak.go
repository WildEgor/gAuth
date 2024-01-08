package keycloak_adapter

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/pkg/errors"
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
	return &KeycloakAdapter{
		gocloak.NewClient(keycloakConfig.API),
		keycloakConfig,
	}
}

func (ka *KeycloakAdapter) loginClient(ctx context.Context) (*gocloak.JWT, error) {
	// TODO: cache JWT token as main keycloak session
	token, err := ka.Client.LoginClient(ctx, ka.keycloakConfig.ClientID, ka.keycloakConfig.ClientSecret, ka.keycloakConfig.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to login the rest client")
	}

	return token, nil
}

// CreateUser Create new user in Keycloak
func (ka *KeycloakAdapter) CreateUser(ctx context.Context, user gocloak.User, password string, role string) (*gocloak.User, error) {
	token, err := ka.loginClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable login to keycloak")
	}

	uID, err := ka.Client.CreateUser(ctx, token.AccessToken, ka.keycloakConfig.Realm, user)
	if err != nil {
		kErr, ok := err.(gocloak.APIError)
		if !ok {
			return nil, errors.Wrap(err, "unable to create the user")
		}
		if kErr.Code == 409 {
			return nil, errors.New("user already exists")
		}
		return nil, errors.Wrap(err, "unable to create the user")
	}

	err = ka.Client.SetPassword(ctx, token.AccessToken, uID, ka.keycloakConfig.Realm, password, false)
	if err != nil {
		return nil, errors.Wrap(err, "unable to set the password for the user")
	}

	var rn = strings.ToLower(role)
	rk, err := ka.Client.GetRealmRole(ctx, token.AccessToken, ka.keycloakConfig.Realm, rn)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to get role by name: '%v'", rn))
	}

	err = ka.Client.AddRealmRoleToUser(ctx, token.AccessToken, ka.keycloakConfig.Realm, uID, []gocloak.Role{
		*rk,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to add a realm role to user")
	}

	uk, err := ka.Client.GetUserByID(ctx, token.AccessToken, ka.keycloakConfig.Realm, uID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get recently created user")
	}

	return uk, nil
}

// Login Auth in Keycloak
func (ka *KeycloakAdapter) Login(user string, pass string) (*JWT, error) {
	token, err := ka.Client.Login(context.Background(), ka.keycloakConfig.ClientID, ka.keycloakConfig.ClientSecret, ka.keycloakConfig.Realm, user, pass)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak login failed")
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

// CheckExists Check user exists in Keycloak
func (ka *KeycloakAdapter) CheckExists(ctx context.Context, user string) (bool, error) {
	token, err := ka.loginClient(ctx)
	if err != nil {
		return false, errors.Wrap(err, "unable login to keycloak")
	}

	search := fmt.Sprintf("%s*", user)
	count, err := ka.Client.GetUserCount(ctx, token.AccessToken, ka.keycloakConfig.Realm, gocloak.GetUsersParams{
		Search: &search,
	})
	if err != nil {
		return false, errors.Wrap(err, "unable to get user count")
	}

	return count > 0, nil
}

// DeleteByEmail Delete user by email
func (ka *KeycloakAdapter) DeleteByEmail(ctx context.Context, email string) error {
	token, err := ka.loginClient(ctx)
	if err != nil {
		return errors.Wrap(err, "unable login to keycloak")
	}

	users, err := ka.Client.GetUsers(ctx, token.AccessToken, ka.keycloakConfig.Realm, gocloak.GetUsersParams{
		Email: &email,
	})
	if err != nil {
		return errors.Wrap(err, "unable to get user by email")
	}

	userId := *users[0].ID

	dErr := ka.Client.DeleteUser(ctx, token.AccessToken, ka.keycloakConfig.Realm, userId)
	if dErr != nil {
		return errors.Wrap(dErr, "unable to delete user")
	}

	logoutErr := ka.Client.LogoutAllSessions(ctx, token.AccessToken, ka.keycloakConfig.Realm, userId)
	if logoutErr != nil {
		return errors.Wrap(logoutErr, "unable to logout user")
	}

	return nil
}

// UpdatePassword Update user password
func (ka *KeycloakAdapter) UpdatePassword(ctx context.Context, email string, newPassword string) error {
	token, err := ka.loginClient(ctx)
	if err != nil {
		return errors.Wrap(err, "unable login to keycloak")
	}

	users, err := ka.Client.GetUsers(ctx, token.AccessToken, ka.keycloakConfig.Realm, gocloak.GetUsersParams{
		Email: &email,
	})
	if err != nil {
		return errors.Wrap(err, "unable to get user by email")
	}

	userId := *users[0].ID

	setErr := ka.Client.SetPassword(ctx, token.AccessToken, userId, ka.keycloakConfig.Realm, newPassword, false)
	if setErr != nil {
		return errors.Wrap(err, "unable to set the password for the user")
	}

	logoutErr := ka.Client.LogoutAllSessions(ctx, token.AccessToken, ka.keycloakConfig.Realm, userId)
	if logoutErr != nil {
		return errors.Wrap(logoutErr, "unable to logout user")
	}

	return nil
}

// UserInfoByToken Get user info by token
func (ka *KeycloakAdapter) UserInfoByToken(ctx context.Context, token string) (*KeycloakUserInfo, error) {
	info, err := ka.Client.GetUserInfo(ctx, token, ka.keycloakConfig.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get user info")
	}

	return &KeycloakUserInfo{
		Id:    *info.Sub,
		Email: *info.Email,
	}, nil
}
