package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	otpHelper "github.com/shu12388y/server/pkg/configs"
	tokenHelper "github.com/shu12388y/server/pkg/configs"
	twilio "github.com/shu12388y/server/pkg/configs"
	database "github.com/shu12388y/server/pkg/database"
	schema "github.com/shu12388y/server/pkg/services/auth/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUpController(c *gin.Context) {
	var resUser schema.User
	var existingUser schema.User

	// Use a reasonable context timeout like 5 seconds instead of 5ms
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	// Bind and validate request body
	if err := c.ShouldBindJSON(&resUser); err != nil {
		log.Println("Error parsing request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	db := database.MongoDBClientConnection()

	// Check if user already exists
	err := db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": resUser.PhoneNumber,
	}).Decode(&existingUser)

	if err == mongo.ErrNoDocuments {
		// User does not exist, proceed with OTP
		otp := fmt.Sprintf("%d", otpHelper.GenerateOTP())
		if otp == "" {
			log.Println("OTP failed to generate")
			return
		}

		_, insertErr := db.Collection("user").InsertOne(ctx, bson.M{
			"phonenumber":  resUser.PhoneNumber,
			"otp":          otp,
			"verified":     false,
			"accesstoken":  "",
			"refreshtoken": "",
		})
		if insertErr != nil {
			log.Println("Error inserting new user:", insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error creating user",
			})
			return
		}

		// Send OTP via WhatsApp
		if sendErr := twilio.SendWhatsAppMessage(resUser.PhoneNumber, otp); sendErr != nil {
			log.Println("Failed to send WhatsApp OTP:", sendErr)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to send OTP",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "OTP sent successfully",
		})
		return
	}

	if err != nil {
		// Some other DB error
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}

	// User already exists
	c.JSON(http.StatusOK, gin.H{
		"message": "User already registered",
	})
}

func VerifyAccount(c *gin.Context) {

	var user schema.User
	var existingUser schema.User

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	c.BindJSON(&user)

	db := database.MongoDBClientConnection()

	// find the user and verify the user
	err := db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&existingUser)
	if err == mongo.ErrNilDocument {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not exists",
		})
		return
	}

	// verify the user and update it
	if user.OTP != existingUser.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong OTP",
		})
		return
	}

	db.Collection("user").FindOneAndUpdate(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}, bson.M{
		"$set": bson.M{
			"verified": true,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func SignInController(c *gin.Context) {

	var user schema.User
	var existingUser schema.User

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	c.BindJSON(&user)

	db := database.MongoDBClientConnection()
	err := db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&existingUser)

	if err == mongo.ErrNilDocument {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "User not exists",
		})
		return
	}
	// check the user verified the user
	if existingUser.Verified == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account is not verified",
		})
		return
	}

	// generate OTP and send
	otp := fmt.Sprintf("%d", otpHelper.GenerateOTP())

	db.Collection("user").FindOneAndUpdate(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	},
		bson.M{
			"$set": bson.M{
				"otp": otp,
			},
		},
	)

	err = twilio.SendWhatsAppMessageLogin(user.PhoneNumber, otp)
	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func SessionController(c *gin.Context) {
	var user schema.User
	var existingUser schema.User

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	c.BindJSON(&user)

	db := database.MongoDBClientConnection()
	err := db.Collection("user").FindOne(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}).Decode(&existingUser)

	if err == mongo.ErrNilDocument {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "User not exists",
		})
		return
	}

	// check the user verified the user
	if existingUser.Verified == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Account is not verified",
		})
		return
	}

	// verify the OTP
	if user.OTP != existingUser.OTP {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Wrong OTP",
		})
	}

	if existingUser.AccessToken != "" && existingUser.RefreshToken != "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"token":   existingUser.RefreshToken,
		})
		return
	}

	// generate the sesssion
	refreshtoken, err := tokenHelper.GenerateRefreshJWT(existingUser.ID.String())
	if err != nil {
		return
	}
	accesstoken, err := tokenHelper.GenerateSecretJWT(existingUser.ID.String())
	if err != nil {
		return
	}

	db.Collection("user").FindOneAndUpdate(ctx, bson.M{
		"phonenumber": user.PhoneNumber,
	}, bson.M{
		"$set": bson.M{
			"accesstoken":  accesstoken,
			"refreshtoken": refreshtoken,
		},
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   accesstoken,
	})

}
