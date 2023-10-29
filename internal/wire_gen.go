// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pkg

import (
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/handlers/health-check"
	"github.com/WildEgor/gAuth/internal/handlers/registration"
	"github.com/WildEgor/gAuth/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from server.go:

func NewServer() (*fiber.App, error) {
	configurator := configs.NewConfigurator()
	appConfig := configs.NewAppConfig(configurator)
	healthCheckHandler := health_check_handler.NewHealthCheckHandler(appConfig)
	registrationHandler := registration_handler.NewRegistrationHandler()
	publicRouter := router.NewPublicRouter(healthCheckHandler, registrationHandler)
	swaggerRouter := router.NewSwaggerRouter()
	mongoDBConfig := configs.NewMongoDBConfig(configurator)
	mongoDBConnection := db.NewMongoDBConnection(mongoDBConfig)
	redisConfig := configs.NewRedisConfig(configurator)
	redisConnection := db.NewRedisDBConnection(redisConfig)
	app := NewApp(appConfig, publicRouter, swaggerRouter, mongoDBConnection, redisConnection)
	return app, nil
}

// server.go:

var ServerSet = wire.NewSet(AppSet)
