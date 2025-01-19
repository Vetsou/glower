package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCart(e *gin.Engine) {
	flowersGroup := e.Group("/cart")

	flowersGroup.POST("/", middleware.VerifyAuthToken, controller.AddFlowerToCart)
}
