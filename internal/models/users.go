package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersModel struct {
	Id           primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	FirstName    string             `json:"first_name,omitempty" bson:"first_name"`
	LastName     string             `json:"last_name,omitempty" bson:"last_name"`
	Email        string             `json:"email" bson:"email"`
	Phone        string             `json:"phone" bson:"phone"`
	Password     string             `json:"password" bson:"password"`
	Verification VerificationModel  `json:"verification,omitempty" bson:"verification"`
	OTP          OTPModel           `json:"otp,omitempty" bson:"otp"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	DeletedAt    time.Time          `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type VerificationModel struct {
	NewPhone     string    `json:"new_phone" bson:"new_phone"`
	NewPhoneCode string    `json:"new_phone_code" bson:"new_phone_code"`
	NewPhoneDate time.Time `json:"new_phone_date,omitempty" bson:"new_phone_date"`
	NewEmail     string    `json:"new_email" bson:"new_email"`
	NewEmailCode string    `json:"new_email_code" bson:"new_email_code"`
	NewEmailDate time.Time `json:"new_email_date,omitempty" bson:"new_email_date"`
}

type OTPModel struct {
	Identity string `json:"identity" bson:"identity"`
	Code     string `json:"code" bson:"code"`
}
