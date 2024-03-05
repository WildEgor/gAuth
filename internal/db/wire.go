package db

import (
	"github.com/WildEgor/gAuth/internal/db/mongo"
	"github.com/WildEgor/gAuth/internal/db/redis"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	mongo.NewMongoConnection,
	redis.NewRedisDBConnection,
)
