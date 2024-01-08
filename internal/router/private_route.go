package router

import (
	cpHandler "github.com/WildEgor/gAuth/internal/handlers/change-password"
	meHandler "github.com/WildEgor/gAuth/internal/handlers/me"
	rtHandler "github.com/WildEgor/gAuth/internal/handlers/refresh"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/WildEgor/gAuth/internal/services"
	"github.com/gofiber/fiber/v2"
)

type PrivateRouter struct {
	cp  *cpHandler.ChangePasswordHandler
	me  *meHandler.MeHandler
	rt  *rtHandler.RefreshHandler
	ur  *repositories.UserRepository
	jwt *services.JWTAuthenticator
}

func NewPrivateRouter(
	cp *cpHandler.ChangePasswordHandler,
	me *meHandler.MeHandler,
	rt *rtHandler.RefreshHandler,
	ur *repositories.UserRepository,
	jwt *services.JWTAuthenticator,
) *PrivateRouter {
	return &PrivateRouter{
		cp:  cp,
		me:  me,
		rt:  rt,
		ur:  ur,
		jwt: jwt,
	}
}

func (r *PrivateRouter) SetupPrivateRouter(app *fiber.App) {
	v1 := app.Group("/api/v1")
	ac := v1.Group("/auth")
	uc := v1.Group("/user")

	am := middlewares.NewAuthMiddleware(middlewares.AuthMiddlewareConfig{
		UR:  r.ur,
		JWT: r.jwt,
	})

	ac.Post("change-password", am, r.cp.Handle)
	ac.Post("refresh", am, r.rt.Handle)
	uc.Get("me", am, r.me.Handle)
}
