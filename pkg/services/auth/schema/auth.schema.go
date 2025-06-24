package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	PhoneNumber  string `bson:"phonenumber" validate:"min=10,max=10"`
	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`
	Verified     bool   `bson:"verified" default:"false"`
	OTP          string `bson:"otp" validate:"min=4,max=4"`
}
