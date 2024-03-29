package router

import (
	"github.com/WildEgor/gAuth/internal/configs"
	hcHandler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	loHandler "github.com/WildEgor/gAuth/internal/handlers/login"
	ltHandler "github.com/WildEgor/gAuth/internal/handlers/logout"
	rcHandler "github.com/WildEgor/gAuth/internal/handlers/reg"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/gofiber/fiber/v2"
)

type PublicRouter struct {
	hc        *hcHandler.HealthCheckHandler
	rc        *rcHandler.RegHandler
	lo        *loHandler.LoginHandler
	lt        *ltHandler.LogoutHandler
	ur        *repositories.UserRepository
	tr        *repositories.TokensRepository
	jwt       *services.JWTAuthenticator
	jwtConfig *configs.JWTConfig
}

func NewPublicRouter(
	hc *hcHandler.HealthCheckHandler,
	rc *rcHandler.RegHandler,
	lo *loHandler.LoginHandler,
	lt *ltHandler.LogoutHandler,
	ur *repositories.UserRepository,
	tr *repositories.TokensRepository,
	jwt *services.JWTAuthenticator,
	jwtConfig *configs.JWTConfig,
) *PublicRouter {
	return &PublicRouter{
		hc:        hc,
		rc:        rc,
		lo:        lo,
		lt:        lt,
		ur:        ur,
		tr:        tr,
		jwt:       jwt,
		jwtConfig: jwtConfig,
	}
}

func (r *PublicRouter) SetupPublicRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// Server endpoint - sanity check that the server is running
	hc := v1.Group("/health")
	hc.Get("check", r.hc.Handle)

	uc := v1.Group("/auth")
	uc.Post("reg", r.rc.Handle)

	lm := middlewares.NewLoginMiddleware(middlewares.LoginMiddlewareConfig{
		UR:        r.ur,
		TR:        r.tr,
		JWT:       r.jwt,
		JWTConfig: r.jwtConfig,
	})
	ac := v1.Group("/auth")

	am := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		UR:  r.ur,
		JWT: r.jwt,
	})
	ac.Post("login", lm, r.lo.Handle)
	ac.Post("logout", am, r.lt.Handle)
}
