package middleware

import (
	"glower/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyAuthToken(c *gin.Context) {
	tokenStr, err := c.Cookie(auth.AccessTokenName)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "error-page.html", gin.H{
			"code":    http.StatusUnauthorized,
			"message": "User is not logged in.",
		})
		c.Abort()
		return
	}

	claims, err := auth.VerifyToken(tokenStr)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "error-page.html", gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid user credentials.",
		})
		c.Abort()
		return
	}

	userData, err := auth.GetUserClaims(claims)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error-page.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error getting user data.",
		})
		c.Abort()
		return
	}

	c.Set("id", userData.Id)
	c.Set("user", userData.User)
	c.Next()
}
