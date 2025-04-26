package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCart(e *gin.Engine) {
	flowersGroup := e.Group("/cart")

	flowersGroup.GET("/", middleware.CreateAuth(true), controller.GetCartItems)
	flowersGroup.POST("/", middleware.CreateAuth(true), controller.AddCartItem)
	flowersGroup.DELETE("/:id", middleware.CreateAuth(true), controller.RemoveCartItem)
}
