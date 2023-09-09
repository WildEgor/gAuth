// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pkg

import (
	"github.com/WildEgor/gAuth/internal/config"
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/handlers/health-check"
	"github.com/WildEgor/gAuth/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

// Injectors from server.go:

func NewServer() (*fiber.App, error) {
	appConfig := config.NewAppConfig()
	healthCheckHandler := health_check_handler.NewHealthCheckHandler(appConfig)
	routerRouter := router.NewRouter(healthCheckHandler)
	mongoDBConfig := config.NewMongoDBConfig()
	mongoDBConnection := db.NewMongoDBConnection(mongoDBConfig)
	redisConfig := config.NewRedisConfig()
	redisConnection := db.NewRedisDBConnection(redisConfig)
	app := NewApp(appConfig, routerRouter, mongoDBConnection, redisConnection)
	return app, nil
}

// server.go:

var ServerSet = wire.NewSet(AppSet)
