package mongo

import (
	"context"
	"github.com/WildEgor/gAuth/internal/configs"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	Client        *mongo.Client
	mongoDbConfig *configs.MongoDBConfig
}

func NewMongoConnection(
	mongoDbConfig *configs.MongoDBConfig,
) *MongoConnection {
	conn := &MongoConnection{
		nil,
		mongoDbConfig,
	}

	conn.Connect()

	return conn
}

func (mc *MongoConnection) Connect() {
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

func (mc *MongoConnection) Disconnect() {
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

func (mc *MongoConnection) AuthDB() *mongo.Database {
	return mc.Client.Database(mc.mongoDbConfig.DbName)
}
