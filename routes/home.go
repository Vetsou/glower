package routes

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
)

func registerHome(e *gin.Engine) {
	homeGroup := e.Group("/")

	homeGroup.GET("/", controller.GetHomePage)
}
