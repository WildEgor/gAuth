package router

import (
	health_check_handler "github.com/WildEgor/gAuth/internal/handlers/health-check"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	hc *health_check_handler.HealthCheckHandler
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

	return nil
}
