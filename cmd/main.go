package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	routes "github.com/shu12388y/server/pkg/services/auth/routes"
)

func main() {
	server := gin.New()

	// health check Router
	server.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	// logger
	server.Use(gin.Logger())

	// routes
	route := server.Group("/api/v1")
	routes.Routes(route)

	// sever port
	server.Run(":5000")
}
