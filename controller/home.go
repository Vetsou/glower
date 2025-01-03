package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
