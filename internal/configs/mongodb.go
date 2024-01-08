package configs

import (
	"github.com/caarlos0/env/v7"
	log "github.com/sirupsen/logrus"
)

type MongoDBConfig struct {
	DbName string `env:"MONGODB_NAME,required"`
	URI    string `env:"MONGODB_URI,required"`
}

func NewMongoDBConfig(c *Configurator) *MongoDBConfig {
	cfg := MongoDBConfig{}

	if err := env.Parse(&cfg); err != nil {
		log.Printf("%+v\n", err)
	}

	return &cfg
}
