package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUser(e *gin.Engine) {
	userGroup := e.Group("/user")

	userGroup.GET("/register", controller.GetRegisterPage)
	userGroup.GET("/login", controller.GetLoginPage)
	userGroup.GET("/profile", middleware.VerifyAuthToken, controller.GetProfilePage)
}
