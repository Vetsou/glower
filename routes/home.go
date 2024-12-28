package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHomeRoute(router *gin.Engine) {
	homeGroup := router.Group("/")

	homeGroup.GET("/", getHome)
}

func getHome(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
