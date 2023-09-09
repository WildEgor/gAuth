package db

import (
	"github.com/WildEgor/gAuth/internal/config"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	Redis *redis.Client
)

type RedisConnection struct {
	redisDbConfig *config.RedisConfig
}

func NewRedisDBConnection(
	redisConfig *config.RedisConfig,
) *RedisConnection {
	conn := &RedisConnection{
		redisConfig,
	}

	conn.connectToRedis()

	return conn
}

func (rc *RedisConnection) connectToRedis() {
	opt, err := redis.ParseURL(rc.redisDbConfig.URI)
	if err != nil {
		panic(err)
	}
	Redis = redis.NewClient(opt)

	_, err = Redis.Ping().Result()
	if err != nil {
		log.Panic("Fail connect to Redis")
		panic(err)
	}

	log.Info("Success connect to Redis")
}

func (rc *RedisConnection) Client() *redis.Client {
	return Redis
}
