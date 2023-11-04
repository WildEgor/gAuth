package repositories

import (
	"time"

	"github.com/WildEgor/gAuth/internal/db"
	"github.com/WildEgor/gAuth/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	DbName = "auth"
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

func (ur *UserRepository) FindByLogin(login string, password string) (*models.UsersModel, error) {
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"phone", bson.D{{"$eq", login}}}},
				bson.D{{"email", bson.D{{"$eq", login}}}},
			},
		},
	}
	var us models.UsersModel
	err := ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).FindOne(nil, filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	hash := []byte(us.Password)
	compPass := []byte(password)
	err = bcrypt.CompareHashAndPassword(hash, compPass)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindById(id string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	var us models.UsersModel
	err := ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).FindOne(nil, filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) Create(nu models.UsersModel) (*models.UsersModel, error) {
	var checkUser models.UsersModel
	checkUser.Email = nu.Email
	checkUser.Phone = nu.Phone

	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"phone", bson.D{{"$eq", nu.Phone}}}},
				bson.D{{"email", bson.D{{"$eq", nu.Email}}}},
			}},
	}
	count, err := ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).CountDocuments(nil, filter)
	if err != nil {
		return nil, errors.Wrap(err, "Mongo error")
	}

	if count > 0 {
		err = errors.New("")
		return nil, errors.Wrap(err, "Email or Phone taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "Generating password hash")
	}
	nu.Password = string(hash)

	us := &models.UsersModel{
		Email:        nu.Email,
		Phone:        nu.Phone,
		Password:     nu.Password,
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		Verification: models.VerificationModel{},
		OTP:          models.OTPModel{},
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err = ur.mongoDbConnection.Instance().Database(DbName).Collection(models.CollectionUsers).InsertOne(nil, us)
	if err != nil {
		return nil, errors.New(`{"mail":"need uniq mail"}`)
	}

	return us, nil
}
