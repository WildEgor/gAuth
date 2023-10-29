package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type MongoDBConfig struct {
	DbName string `json:"MONGODB_NAME"`
	URI    string `env:"MONGODB_URI"`
}

func NewMongoDBConfig(c *Configurator) *MongoDBConfig {
	cfg := MongoDBConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("[MongoDBConfig] %+v\n", err)
	}

	return &cfg
}
