package router

import (
	hcHandler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	loHandler "github.com/WildEgor/gAuth/internal/handlers/login"
	rcHandler "github.com/WildEgor/gAuth/internal/handlers/reg"
	"github.com/WildEgor/gAuth/internal/middlewares"
	"github.com/WildEgor/gAuth/internal/repositories"
	"github.com/gofiber/fiber/v2"
)

type PublicRouter struct {
	hc *hcHandler.HealthCheckHandler
	rc *rcHandler.RegHandler
	lo *loHandler.LoginHandler
	ur *repositories.UserRepository
}

func NewPublicRouter(
	hc *hcHandler.HealthCheckHandler,
	rc *rcHandler.RegHandler,
	lo *loHandler.LoginHandler,
	ur *repositories.UserRepository,
) *PublicRouter {
	return &PublicRouter{
		hc: hc,
		rc: rc,
		lo: lo,
		ur: ur,
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
		UserRepo: r.ur,
	})
	ac := v1.Group("/auth")
	ac.Post("login", lm, r.lo.Handle)
}
