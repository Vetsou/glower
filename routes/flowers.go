package routes

import (
	"glower/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterFlowersRoute(router *gin.Engine) {
	flowersGroup := router.Group("/flowers")

	flowersGroup.GET("/", getFlowers)
}

func getFlowers(c *gin.Context) {
	flowers, err := controller.GetFlowersStock()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to load flowers. Please try again later.",
		})
		return
	}

	c.HTML(http.StatusOK, "shop-stock.html", gin.H{
		"flowers": flowers,
	})
}
