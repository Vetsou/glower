package middleware

import (
	"glower/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyAuthToken(c *gin.Context) {
	tokenStr, err := c.Cookie("access-token")
	if err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": "Access token is missing.",
		})
		c.Abort()
		return
	}

	claims, err := utils.VerifyToken(tokenStr)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"error": "Invalid token.",
		})
		c.Abort()
		return
	}

	c.Set("user", (*claims)["user"])

	c.Next()
}
