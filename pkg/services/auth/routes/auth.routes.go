package services

import (
	"github.com/gin-gonic/gin"
	controller "github.com/shu12388y/server/pkg/services/auth/controller"
)

func Routes(r *gin.RouterGroup) {
	r.POST("/signup", controller.SignUpController)
	r.POST("/verify", controller.VerifyAccount)
	r.POST("/signin", controller.SignInController)
	r.POST("/session", controller.SessionController)
}
