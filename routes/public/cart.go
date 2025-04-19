package public

import (
	"glower/controller"
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCart(e *gin.Engine) {
	flowersGroup := e.Group("/cart")

	flowersGroup.GET("/", middleware.CreateAuthMiddleware(true), controller.GetCartItems)
	flowersGroup.POST("/", middleware.CreateAuthMiddleware(true), controller.AddCartItem)
	flowersGroup.DELETE("/:id", middleware.CreateAuthMiddleware(true), controller.RemoveCartItem)
}
