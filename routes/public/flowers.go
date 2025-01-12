package public

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
)

func RegisterFlowers(e *gin.Engine) {
	flowersGroup := e.Group("/flowers")

	flowersGroup.GET("/", controller.GetFlowers)
	flowersGroup.POST("/", controller.AddFlower)
	flowersGroup.DELETE("/:id", controller.RemoveFlower)
}
