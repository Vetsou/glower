package routes

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
)

func registerFlowers(e *gin.Engine) {
	flowersGroup := e.Group("/flowers")

	flowersGroup.GET("/", controller.GetFlowers)
	flowersGroup.POST("/", controller.AddFlower)
	flowersGroup.DELETE("/:id", controller.RemoveFlower)
}
