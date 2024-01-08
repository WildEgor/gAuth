package repositories

import (
	"github.com/google/wire"
)

var RepositoriesSet = wire.NewSet(
	NewUserRepository,
	NewTokensRepository,
)
