package routes

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func registerUser(e *gin.Engine) {
	userGroup := e.Group("/user")

	userGroup.GET("/register", controller.GetRegisterPage)
	userGroup.GET("/login", controller.GetLoginPage)
	userGroup.GET("/profile", middleware.VerifyAuthToken, controller.GetProfilePage)
}
