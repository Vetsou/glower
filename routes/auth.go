package routes

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func registerAuth(e *gin.Engine) {
	authGroup := e.Group("/auth")

	authGroup.POST("/signup", controller.RegisterUser)
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/logout", middleware.VerifyAuthToken, controller.Logout)
}
