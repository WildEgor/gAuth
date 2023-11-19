package db

import (
	"context"
	"github.com/WildEgor/gAuth/internal/configs"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnection struct {
	Client        *mongo.Client
	mongoDbConfig *configs.MongoDBConfig
}

func NewMongoDBConnection(
	mongoDbConfig *configs.MongoDBConfig,
) *MongoDBConnection {
	conn := &MongoDBConnection{
		nil,
		mongoDbConfig,
	}

	defer conn.Disconnect()

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

	log.Info("Success connect to MongoDB")

	mc.Client = client
}

func (mc *MongoDBConnection) Disconnect() {
	if mc.Client == nil {
		return
	}

	err := mc.Client.Disconnect(context.TODO())
	if err != nil {
		log.Panic("Fail disconnect Mongo", err)
		panic(err)
	}

	log.Info("Connection to MongoDB closed.")
}

func (mc *MongoDBConnection) Instance() *mongo.Client {
	return mc.Client
}

func (mc *MongoDBConnection) DbName() string {
	return mc.mongoDbConfig.DbName
}
