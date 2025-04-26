package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUser(e *gin.Engine) {
	userGroup := e.Group("/user")

	userGroup.GET("/register", middleware.CreateAuth(false), controller.GetRegisterPage)
	userGroup.GET("/login", middleware.CreateAuth(false), controller.GetLoginPage)
	userGroup.GET("/profile", middleware.CreateAuth(true), controller.GetProfilePage)
}
