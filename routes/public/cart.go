package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCart(e *gin.Engine) {
	flowersGroup := e.Group("/cart")

	flowersGroup.GET("/", middleware.VerifyAuthToken, controller.GetCartItems)
	flowersGroup.POST("/", middleware.VerifyAuthToken, controller.AddCartItem)
	flowersGroup.DELETE("/:id", middleware.VerifyAuthToken, controller.RemoveCartItem)
}
