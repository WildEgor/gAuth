package pkg

import (
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/WildEgor/gAuth/internal/proto"
	"os"

	"github.com/WildEgor/gAuth/internal/adapters"
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/wire"

	error_handler "github.com/WildEgor/gAuth/internal/errors"
	log "github.com/sirupsen/logrus"
)

var AppSet = wire.NewSet(
	NewApp,
	proto.GRPCSet,
	configs.ConfigsSet,
	router.RouterSet,
	db.DbSet,
	adapters.AdaptersSet,
)

func NewApp(
	appConfig *configs.AppConfig,
	pbRouter *router.PublicRouter,
	prRouter *router.PrivateRouter,
	swaggerRouter *router.SwaggerRouter,
	mongo *db.MongoDBConnection,
	redis *db.RedisConnection,
) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: error_handler.ErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	app.Use(recover.New())

	// Set logging settings
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	if !appConfig.IsProduction() {
		// HINT: some extra setting
		log.SetLevel(log.DebugLevel)
		log.SetLevel(log.ErrorLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	prRouter.SetupPrivateRouter(app)
	pbRouter.SetupPublicRouter(app)
	swaggerRouter.SetupSwaggerRouter(app)

	mongo.Connect()
	redis.Connect()

	return app
}
