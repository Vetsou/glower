package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUser(e *gin.Engine) {
	userGroup := e.Group("/user")

	userGroup.GET("/register",
		middleware.CreateAuth(false),
		controller.CreateRegisterPage())

	userGroup.GET("/login",
		middleware.CreateAuth(false),
		controller.CreateLoginPage())

	userGroup.GET("/profile",
		middleware.CreateAuth(true),
		controller.CreateProfilePage())
}
