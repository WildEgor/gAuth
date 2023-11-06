package router

import (
	kcAdapter "github.com/WildEgor/gAuth/internal/adapters/keycloak"
	"github.com/WildEgor/gAuth/internal/configs"
	cpHandler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type PrivateRouter struct {
	cp *cpHandler.ChangePasswordHandler
	ka *kcAdapter.KeycloakAdapter
	kc *configs.KeycloakConfig
	ur *repositories.UserRepository
}

func NewPrivateRouter(
	cp *cpHandler.ChangePasswordHandler,
	ka *kcAdapter.KeycloakAdapter,
	kc *configs.KeycloakConfig,
	ur *repositories.UserRepository,
) *PrivateRouter {
	return &PrivateRouter{
		cp: cp,
		ka: ka,
		kc: kc,
		ur: ur,
	}
}

func (r *PrivateRouter) SetupPrivateRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")

	authMiddleware := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		KeycloakAdapter: r.ka,
		KeycloakConfig:  r.kc,
		UserRepo:        r.ur,
	})
	authController := v1.Group("/auth")
	authController.Post("change-password", authMiddleware, r.cp.Handle)
}
