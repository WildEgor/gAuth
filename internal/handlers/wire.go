package handlers

import (
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	me_handler "github.com/WildEgor/gAuth/internal/handlers/me"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	health_check_handler.NewHealthCheckHandler,
	me_handler.NewMeHandler,
)
