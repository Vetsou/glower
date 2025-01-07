package initializers

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterServerRoutes(e *gin.Engine) {
	homeGroup := e.Group("/")
	{
		homeGroup.GET("/", controller.GetHomePage)
	}

	flowersGroup := e.Group("/flowers")
	{
		flowersGroup.GET("/", controller.GetFlowers)
		flowersGroup.POST("/", controller.AddFlower)
		flowersGroup.DELETE("/:id", controller.RemoveFlower)
	}

	authGroup := e.Group("/auth")
	{
		authGroup.POST("/signup", controller.RegisterUser)
		authGroup.POST("/login", controller.Login)
		authGroup.POST("/logout", middleware.VerifyAuthToken, controller.Logout)
	}

	userGroup := e.Group("/user")
	{
		userGroup.GET("/register", controller.GetRegisterPage)
		userGroup.GET("/login", controller.GetLoginPage)
		userGroup.GET("/profile", middleware.VerifyAuthToken, controller.GetProfilePage)
	}
}
