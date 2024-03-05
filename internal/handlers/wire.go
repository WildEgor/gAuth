package handlers

import (
	change_password_handler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	login_handler "github.com/WildEgor/gAuth/internal/handlers/login"
	logout_handler "github.com/WildEgor/gAuth/internal/handlers/logout"
	me_handler "github.com/WildEgor/gAuth/internal/handlers/me"
	otp_generate_handler "github.com/WildEgor/gAuth/internal/handlers/otp-generate"
	otp_login_handler "github.com/WildEgor/gAuth/internal/handlers/otp-login"
	refresh_handler "github.com/WildEgor/gAuth/internal/handlers/refresh"
	registration_handler "github.com/WildEgor/gAuth/internal/handlers/reg"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/google/wire"
)

var HandlersSet = wire.NewSet(
	repositories.RepositoriesSet,
	services.ServicesSet,
	health_check_handler.NewHealthCheckHandler,
	me_handler.NewMeHandler,
	registration_handler.NewRegHandler,
	refresh_handler.NewRefreshHandler,
	login_handler.NewLoginHandler,
	logout_handler.NewLogoutHandler,
	change_password_handler.NewChangePasswordHandler,
	otp_generate_handler.NewOTPGenHandler,
	otp_login_handler.NewOTPLoginHandler,
)
