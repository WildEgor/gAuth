package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersModel struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Email     string             `json:"email,omitempty" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt time.Time          `json:"deleted_at" bson:"deleted_at"`
}
