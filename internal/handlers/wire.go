package handlers

import (
	centrifuge_validate_token_handler "github.com/WildEgor/gAuth/internal/handlers/centrifuge-validate-token"
	change_password_handler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	login_handler "github.com/WildEgor/gAuth/internal/handlers/login"
	me_handler "github.com/WildEgor/gAuth/internal/handlers/me"
	registration_handler "github.com/WildEgor/gAuth/internal/handlers/registration"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	health_check_handler.NewHealthCheckHandler,
	me_handler.NewMeHandler,
	registration_handler.NewRegistrationHandler,
	login_handler.NewLoginHandler,
	change_password_handler.NewChangePasswordHandler,
	centrifuge_validate_token_handler.NewCentrifugeValidateToken,
)
