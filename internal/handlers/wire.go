package handlers

import (
	change_password_handler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	login_handler "github.com/WildEgor/gAuth/internal/handlers/login"
	me_handler "github.com/WildEgor/gAuth/internal/handlers/me"
	refresh_handler "github.com/WildEgor/gAuth/internal/handlers/refresh"
	registration_handler "github.com/WildEgor/gAuth/internal/handlers/reg"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	health_check_handler.NewHealthCheckHandler,
	me_handler.NewMeHandler,
	registration_handler.NewRegHandler,
	refresh_handler.NewRefreshHandler,
	login_handler.NewLoginHandler,
	change_password_handler.NewChangePasswordHandler,
)
