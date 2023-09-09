package repositories

import (
	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	mongoDbConnection *db.MongoDBConnection
}

func NewUserRepository(
	mongoDbConnection *db.MongoDBConnection,
) *UserRepository {
	return &UserRepository{
		mongoDbConnection,
	}
}

func (ur *UserRepository) FindById(id string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	res := ur.mongoDbConnection.Session().Database("q-auth").Collection(models.CollectionUsers).FindOne(nil, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var us *models.UsersModel
	res.Decode(us)
	return us, nil
}
