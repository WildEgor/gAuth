package router

import (
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	hcHandler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	loHandler "github.com/WildEgor/gAuth/internal/handlers/login"
	rcHandler "github.com/WildEgor/gAuth/internal/handlers/registration"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type PublicRouter struct {
	hc *hcHandler.HealthCheckHandler
	rc *rcHandler.RegistrationHandler
	lo *loHandler.LoginHandler
	ka *kcAdapter.KeycloakAdapter
	ur *repositories.UserRepository
}

func NewPublicRouter(
	hc *hcHandler.HealthCheckHandler,
	rc *rcHandler.RegistrationHandler,
	lo *loHandler.LoginHandler,
	ka *kcAdapter.KeycloakAdapter,
	ur *repositories.UserRepository,
) *PublicRouter {
	return &PublicRouter{
		hc: hc,
		rc: rc,
		lo: lo,
		ka: ka,
		ur: ur,
	}
}

func (r *PublicRouter) SetupPublicRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")

	// Server endpoint - sanity check that the server is running
	hcController := v1.Group("/health")
	hcController.Get("check", r.hc.Handle)

	userController := v1.Group("/user")
	userController.Post("reg", r.rc.Handle)

	loginMiddleware := middlewares.NewLoginMiddleware(middlewares.LoginMiddlewareConfig{
		KeycloakAdapter: r.ka,
		UserRepo:        r.ur,
	})
	authController := v1.Group("/auth")
	authController.Post("login", loginMiddleware, r.lo.Handle)
}
