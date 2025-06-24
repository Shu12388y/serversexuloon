package services

import (
	"context"
	"github.com/gin-gonic/gin"
	database "github.com/shu12388y/server/pkg/database"
	schema "github.com/shu12388y/server/pkg/services/auth/schema"
	otpHelper "github.com/shu12388y/server/pkg/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func SignUpController(c *gin.Context) {

	// userSchema
	var user schema.User
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the json data
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Not able to parse the info",
		})
		log.Println("BindJSON Error", err)
		return
	}

	db := database.MongoDBClientConnection()

	err = db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&user)
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User already exists",
		})
		return
	}

	// Generate the OTP
	generateOTP := otpHelper.GenerateOTP()


	// Insert the user 
	res, err := db.Collection("user").InsertOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
		"verified":    false,
		"otp": generateOTP,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server error",
		})
		return
	}
	if res.InsertedID != nil {
		// send the OTP
		log.Println("User Inserted")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})

}




func SignInController(c *gin.Context){

	// User schema
	var user schema.User


	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second);
	defer cancel()





	



}
