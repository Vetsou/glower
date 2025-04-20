package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUser(e *gin.Engine) {
	userGroup := e.Group("/user")

	userGroup.GET("/register", middleware.CreateAuthMiddleware(false), controller.GetRegisterPage)
	userGroup.GET("/login", middleware.CreateAuthMiddleware(false), controller.GetLoginPage)
	userGroup.GET("/profile", middleware.CreateAuthMiddleware(true), controller.GetProfilePage)
}
