package db

import (
	"context"

	"github.com/WildEgor/gAuth/internal/config"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Session *mongo.Client
)

type MongoDBConnection struct {
	mongoDbConfig *config.MongoDBConfig
}

func NewMongoDBConnection(
	mongoDbConfig *config.MongoDBConfig,
) *MongoDBConnection {
	conn := &MongoDBConnection{
		mongoDbConfig,
	}

	defer conn.disconnect()

	return conn
}

func (mc *MongoDBConnection) Connect() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mc.mongoDbConfig.URI))
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Panic("Fail connect to Mongo", err)
	}

	Session = client

	log.Info("Success connect to MongoDB")
}

func (mc *MongoDBConnection) disconnect() {
	if Session == nil {
		return
	}

	err := Session.Disconnect(context.TODO())
	if err != nil {
		log.Panic("Fail disconnect Mongo", err)
		panic(err)
	}

	log.Info("Connection to MongoDB closed.")
}

func (mc *MongoDBConnection) Session() *mongo.Client {
	return Session
}

func (mc *MongoDBConnection) DbName() string {
	return mc.mongoDbConfig.DbName
}
