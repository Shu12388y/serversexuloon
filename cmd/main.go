package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	routes "github.com/shu12388y/server/pkg/services/auth/routes"
	"net/http"
	"time"
)

func main() {
	server := gin.New()

	// configuring middlewares
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // or ["*"] for all
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
