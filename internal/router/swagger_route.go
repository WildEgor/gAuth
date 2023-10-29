package router

import (
	"github.com/gofiber/fiber/v2"

	swagger "github.com/gofiber/swagger"
)

type SwaggerRouter struct {
}

func NewSwaggerRouter() *SwaggerRouter {
	return &SwaggerRouter{}
}

// SetupSwaggerRouter func for describe group of API Docs routes.
func (sr *SwaggerRouter) SetupSwaggerRouter(app *fiber.App) error {
	// Create routes group.
	route := app.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault)

	return nil
}
