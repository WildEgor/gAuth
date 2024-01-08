package router

import (
	cpHandler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	rtHandler "github.com/WildEgor/gAuth/internal/handlers/refresh"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type PrivateRouter struct {
	cp *cpHandler.ChangePasswordHandler
	rt *rtHandler.RefreshHandler
	ur *repositories.UserRepository
}

func NewPrivateRouter(
	cp *cpHandler.ChangePasswordHandler,
	rt *rtHandler.RefreshHandler,
	ur *repositories.UserRepository,
) *PrivateRouter {
	return &PrivateRouter{
		cp: cp,
		rt: rt,
		ur: ur,
	}
}

func (r *PrivateRouter) SetupPrivateRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")
	ac := v1.Group("/auth")

	am := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		UserRepo: r.ur,
	})

	ac.Post("change-password", am, r.cp.Handle)
	ac.Post("refresh", am, r.rt.Handle)
}
