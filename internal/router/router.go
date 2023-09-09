package router

import (
	keycloak_adapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	hc              *health_check_handler.HealthCheckHandler
	keycloakAdapter *keycloak_adapter.KeycloakAdapter
}

func NewRouter(hc *health_check_handler.HealthCheckHandler) *Router {
	return &Router{
		hc: hc,
	}
}

func (r *Router) Setup(app *fiber.App) error {
	v1 := app.Group("/api/v1")

	// Server endpoint - sanity check that the server is running
	hcController := v1.Group("/health")
	hcController.Get("/check", r.hc.Handle)

	authMiddleware := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		KeycloakAdapter: r.keycloakAdapter,
	})

	authController := v1.Group("/auth")
	authController.Post("/login", authMiddleware)

	return nil
}
