package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	schema "github.com/shu12388y/server/pkg/services/auth/schema"
)

func SignUpController(c *gin.Context) {

	// userSchema
	var user schema.User

	// get the json data
	err := c.BindJSON(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Not able to parse the info",
		})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"user": user,
		},
	})

}
