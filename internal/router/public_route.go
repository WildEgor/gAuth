package router

import (
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	hcHandler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	rcHandler "github.com/WildEgor/gAuth/internal/handlers/registration"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type PublicRouter struct {
	hc              *hcHandler.HealthCheckHandler
	rc              *rcHandler.RegistrationHandler
	keycloakAdapter *kcAdapter.KeycloakAdapter
}

func NewPublicRouter(
	hc *hcHandler.HealthCheckHandler,
	rc *rcHandler.RegistrationHandler,
) *PublicRouter {
	return &PublicRouter{
		hc: hc,
		rc: rc,
	}
}

func (r *PublicRouter) SetupPublicRouter(app *fiber.App) error {
	v1 := app.Group("/api/v1")

	// Server endpoint - sanity check that the server is running
	hcController := v1.Group("/health")
	hcController.Get("/check", r.hc.Handle)

	authMiddleware := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		KeycloakAdapter: r.keycloakAdapter,
	})

	authController := v1.Group("/auth")
	authController.Post("/login", authMiddleware)

	userController := v1.Group("/user")
	userController.Post("/reg", r.rc.Handle)

	return nil
}
