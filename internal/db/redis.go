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

	defer conn.disconnect()

	return conn
}

func (rc *RedisConnection) Connect() {
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

func (mc *RedisConnection) disconnect() {
	if Redis == nil {
		return
	}

	err := Redis.Close()
	if err != nil {
		log.Panic("Fail disconnect Redis", err)
		panic(err)
	}

	log.Info("Connection to Redis closed.")
}

func (rc *RedisConnection) Client() *redis.Client {
	return Redis
}
