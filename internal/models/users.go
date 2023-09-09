package models

import (
	"time"

	"github.com/golang-jwt/jwt"
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

type WebUserDTO struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"admin"`
}

type JwtUserTokenDTO struct {
	WebUserDTO
	jwt.StandardClaims
}

type JwtUserTokenResponseDTO struct {
	Token string `json:"token"`
}

type LoginUserRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewUserRequestDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
