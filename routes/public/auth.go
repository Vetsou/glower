package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuth(e *gin.Engine) {
	authGroup := e.Group("/auth")

	authGroup.POST("/signup", controller.RegisterUser)
	authGroup.POST("/login", controller.Login)
	authGroup.POST("/logout", middleware.CreateAuthMiddleware(true), controller.Logout)
}
