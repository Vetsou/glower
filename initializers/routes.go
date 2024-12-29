package initializers

import (
	"glower/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterServerRoutes(e *gin.Engine) {
	homeGroup := e.Group("/")
	{
		homeGroup.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.html", nil) })
	}

	flowersGroup := e.Group("/flowers")
	{
		flowersGroup.GET("/", controller.GetFlowers)
		flowersGroup.POST("/", controller.AddFlower)
		flowersGroup.DELETE("/:id", controller.RemoveFlower)
	}
}
