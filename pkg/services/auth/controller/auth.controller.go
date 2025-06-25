package services

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	otpHelper "github.com/shu12388y/server/pkg/configs"
	database "github.com/shu12388y/server/pkg/database"
	schema "github.com/shu12388y/server/pkg/services/auth/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUpController(c *gin.Context) {
	// userSchema
	var user schema.User

	// result user body
	var resUser schema.User

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
	}).Decode(&resUser)
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
		"otp":         generateOTP,
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

func SignInOTPController(c *gin.Context) {

	// User schema
	var user schema.User

	// result user
	var resUser schema.User

	// context for handling async code
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the user
	err := c.BindJSON(&user)
	if err != nil {
		log.Println("User empty request error")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Empty request",
		})
		return
	}

	// find the user in the user exist or not
	db := database.MongoDBClientConnection()
	err = db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&resUser)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not exists",
		})
	}

	// check the user verified or not
	if resUser.Verified == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Verify the account first",
		})
		return
	}

	// if user is verified generate the OTP and send the OTP use the queue and
	var generateOTP = otpHelper.GenerateOTP()

	// update the OTP in database
	result := db.Collection("user").FindOneAndUpdate(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}, bson.M{
		"$set": bson.M{
			"otp": generateOTP,
		},
	})
	if result.Err() != nil {
		log.Println("OTP updating error", result.Err())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "succcess",
	})

}

func VerifyController(c *gin.Context) {

	// define the user
	var user schema.User

	// define the resulted user
	var resUser schema.User

	// get the json data
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Not able to parse the info",
		})
		log.Println("BindJSON Error", err)
		return
	}

	// define context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	// find the user
	db := database.MongoDBClientConnection()
	err = db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&resUser)

	if err == mongo.ErrNilDocument {
		log.Println("User not exist send in verify request", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	if user.OTP != resUser.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong OTP",
		})
		return
	}

	// update the user
	result := db.Collection("user").FindOneAndUpdate(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}, bson.M{
		"$set": bson.M{
			"verified": true,
		},
	})
	if result.Err() != nil {
		log.Println("Failed to verify and update in the DB", result.Err())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func SignInController(c *gin.Context) {

	// user schema
	var user schema.User

	// res user
	var resUser schema.User

	// creating context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	// get the json body
	c.BindJSON(&user)

	// find the user
	db := database.MongoDBClientConnection()
	err := db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&resUser)
	if err == mongo.ErrNilDocument {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "User not exists",
		})
	}

	// check the OTP
	if user.OTP != resUser.OTP {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Wrong OTP",
		})
		return
	}

	// if OTP is verified create the Refresh Token and Access Token

}
