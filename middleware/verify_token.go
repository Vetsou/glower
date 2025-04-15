package middleware

import (
	"glower/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func renderResponse(c *gin.Context, code int, msg string) {
	isHTMX := c.Request.Header.Get("HX-Request") == "true"

	if isHTMX {
		c.HTML(code, "error-alert.html", gin.H{"errorMessage": msg})
	} else {
		c.HTML(code, "error-page.html", gin.H{
			"code":    code,
			"message": msg,
		})
	}
	c.Abort()
}

func VerifyAuthToken(c *gin.Context) {
	tokenStr, err := c.Cookie(auth.AccessTokenName)
	if err != nil {
		renderResponse(c, http.StatusUnauthorized, "User is not logged in.")
		return
	}

	claims, err := auth.VerifyToken(tokenStr)
	if err != nil {
		renderResponse(c, http.StatusUnauthorized, "Invalid user credentials.")
		return
	}

	userData, err := auth.GetUserClaims(claims)
	if err != nil {
		renderResponse(c, http.StatusInternalServerError, "Error getting user data.")
		return
	}

	c.Set("id", userData.Id)
	c.Set("user", userData.User)
	c.Next()
}
