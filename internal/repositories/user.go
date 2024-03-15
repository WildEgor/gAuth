package repositories

import (
	"context"
	mongo2 "github.com/WildEgor/gAuth/internal/db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/WildEgor/gAuth/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(
	db *mongo2.MongoConnection,
) *UserRepository {

	coll := db.AuthDB().Collection(models.CollectionUsers)

	return &UserRepository{
		coll,
	}
}

func (ur *UserRepository) FindByPhone(phone string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "phone", Value: phone}}

	var us models.UsersModel
	err := ur.coll.FindOne(context.TODO(), filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByEmail(email string) (*models.UsersModel, error) {
	filter := bson.D{{Key: "email", Value: email}}

	var us models.UsersModel
	err := ur.coll.FindOne(context.TODO(), filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByLogin(login string) (*models.UsersModel, error) {
	filter := bson.D{
		{"$or",
			bson.A{
				bson.D{{"phone", bson.D{{"$eq", login}}}},
				bson.D{{"email", bson.D{{"$eq", login}}}},
			},
		},
	}
	var us models.UsersModel
	err := ur.coll.FindOne(context.TODO(), filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) FindByIds(ids []string) (*[]models.UsersModel, error) {

	objectIds := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objectIds[i], _ = primitive.ObjectIDFromHex(id)
	}

	filter := bson.D{{"_id", bson.D{{"$in", objectIds}}}}

	cursor, err := ur.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var users []models.UsersModel
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return &users, nil
}

func (ur *UserRepository) CountAll() (int64, error) {
	count, err := ur.coll.CountDocuments(context.TODO(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "Mongo error")
	}

	return count, nil
}

func (ur *UserRepository) FindById(id string) (*models.UsersModel, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{Key: "_id", Value: oid}}

	var us models.UsersModel
	err := ur.coll.FindOne(context.TODO(), filter).Decode(&us)
	if err != nil {
		return nil, err
	}

	return &us, nil
}

func (ur *UserRepository) Create(nu *models.UsersModel) (*models.UsersModel, error) {
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
	count, err := ur.coll.CountDocuments(context.TODO(), filter)
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
		Status:       models.ActiveStatus,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	insertResult, err := ur.coll.InsertOne(context.TODO(), us)
	if err != nil {
		return nil, errors.New(`{"mail":"need uniq mail"}`)
	}

	us.Id = insertResult.InsertedID.(primitive.ObjectID)

	return us, nil
}

func (ur *UserRepository) UpdatePassword(nu models.UsersModel) error {
	nu.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{"$set",
			bson.D{
				{"password", nu.Password},
				{"updated_at", nu.UpdatedAt},
			},
		},
	}

	_, err := ur.coll.UpdateByID(context.TODO(), nu.Id, update)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) UpdateOTP(id primitive.ObjectID, otp models.OTPModel) error {
	update := bson.D{
		{"$set",
			bson.D{
				{"otp", otp},
				{"updated_at", time.Now().UTC()},
			},
		},
	}

	_, err := ur.coll.UpdateByID(context.TODO(), id, update)
	if err != nil {
		return err
	}

	return nil
}
