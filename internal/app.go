package pkg

import (
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/db/mongo"
	"github.com/WildEgor/gAuth/internal/db/redis"
	"github.com/WildEgor/gAuth/internal/proto"
	"github.com/WildEgor/gAuth/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/wire"
	"os"

	errorhandler "github.com/WildEgor/gAuth/internal/errors"
	log "github.com/sirupsen/logrus"
)

var AppSet = wire.NewSet(
	NewApp,
	configs.ConfigsSet,
	router.RouterSet,
	proto.RPCSet,
	db.DbSet,
)

type Server struct {
	App       *fiber.App
	AppConfig *configs.AppConfig
	GRPC      *proto.GRPCServer
	Mongo     *mongo.MongoConnection
	Redis     *redis.RedisConnection
}

func NewApp(
	appConfig *configs.AppConfig,
	pbRouter *router.PublicRouter,
	prRouter *router.PrivateRouter,
	swaggerRouter *router.SwaggerRouter,
	server *proto.GRPCServer,
	mongo *mongo.MongoConnection,
	redis *redis.RedisConnection,
) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.ErrorHandler,
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

	return &Server{
		App:       app,
		AppConfig: appConfig,
		Redis:     redis,
		Mongo:     mongo,
		GRPC:      server,
	}
}
