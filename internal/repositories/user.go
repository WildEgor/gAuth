package repositories

import (
	"time"

	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	DbName = "test"
)

type UserRepository struct {
	mongoDbConnection *db.MongoDBConnection
}

func NewUserRepository(
	mongoDbConnection *db.MongoDBConnection,
) *UserRepository {
	DbName = mongoDbConnection.DbName()

	return &UserRepository{
		mongoDbConnection,
	}
}

func (ur *UserRepository) FindById(id string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	res := ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).FindOne(nil, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var us *models.UsersModel
	res.Decode(us)
	return us, nil
}

func (ur *UserRepository) FindCount(user models.UsersModel) (int64, error) {
	count, err := ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).CountDocuments(nil, user)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (ur *UserRepository) Create(nu models.UsersModel) (*models.UsersModel, error) {
	var checkUser models.UsersModel
	checkUser.Email = nu.Email
	result, err := ur.FindCount(checkUser)
	if result > 0 {
		return &nu, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "Generating password hash")
	}
	nu.Password = string(hash)

	us := &models.UsersModel{
		Id:        primitive.NewObjectID(),
		Email:     nu.Email,
		Password:  nu.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).InsertOne(nil, us)
	if err != nil {
		return nil, errors.New(`{"mail":"need uniq mail"}`)
	}

	return us, nil
}
