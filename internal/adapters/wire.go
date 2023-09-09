package adapters

import (
	keycloak_adapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/google/wire"
)

var AdaptersSet = wire.NewSet(
	keycloak_adapter.NewKeycloakAdapter,
)
