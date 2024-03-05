package redis

import (
	"github.com/WildEgor/gAuth/internal/configs"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type RedisConnection struct {
	Client        *redis.Client
	redisDbConfig *configs.RedisConfig
}

func NewRedisDBConnection(
	redisConfig *configs.RedisConfig,
) *RedisConnection {
	conn := &RedisConnection{
		nil,
		redisConfig,
	}

	conn.Connect()

	return conn
}

func (rc *RedisConnection) Connect() {
	opt, err := redis.ParseURL(rc.redisDbConfig.URI)
	if err != nil {
		log.Panic("Fail parse URL ", err)
		panic(err)
	}

	rc.Client = redis.NewClient(opt)

	if _, err := rc.Client.Ping().Result(); err != nil {
		log.Panic("Fail connect to Redis ", err)
		panic(err)
	}

	log.Info("Success connect to Redis")
}

func (rc *RedisConnection) Disconnect() {
	if rc.Client == nil {
		return
	}

	if err := rc.Client.Close(); err != nil {
		log.Panic("Fail disconnect Redis", err)
		panic(err)
	}

	log.Info("Connection to Redis closed.")
}

func (rc *RedisConnection) Instance() *redis.Client {
	return rc.Client
}
