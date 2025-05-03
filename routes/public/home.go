package public

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
)

func RegisterHome(e *gin.Engine) {
	homeGroup := e.Group("/")

	homeGroup.GET("/", controller.CreateGetHomePage())
}
