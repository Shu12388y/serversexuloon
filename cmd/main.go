package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	queue "github.com/shu12388y/server/pkg/queues"
)

func main() {
	server := gin.New()

	// health check Router
	server.GET("/health", func(ctx *gin.Context) {
		// testing queues
		messageQueue := make(chan string, 10)
		job := queue.ProducerJob{
			Payload: map[string]string{
				"data": "9905575178",
			},
			Event: "OTP-VERIFY",
		}
		go queue.Producers(job, messageQueue)

		queue.Consumers(messageQueue)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ok",
		})
	})

	// logger
	server.Use(gin.Logger())

	// sever port
	server.Run(":5000")
}
