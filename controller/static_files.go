package controller

import (
	res "glower/resources"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateServeFavicon() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := res.AssetsFS.ReadFile("assets/favicon.ico")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Data(http.StatusOK, "image/x-icon", data)
	}
}
