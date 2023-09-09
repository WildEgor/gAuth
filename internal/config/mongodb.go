package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type MongoDBConfig struct {
	URI string `env:"MONGODB_URI"`
}

func NewMongoDBConfig() *MongoDBConfig {
	cfg := MongoDBConfig{}

	if err := godotenv.Load(".env"); err == nil {
		if err := env.Parse(&cfg); err != nil {
			log.Printf("[MongoDBConfig] %+v\n", err)
		}
	}

	return &cfg
}
