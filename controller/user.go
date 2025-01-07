package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user-register.html", nil)
}

func GetLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user-login.html", nil)
}

func GetProfilePage(c *gin.Context) {
	val, exists := c.Get("user")

	if !exists {
		c.HTML(http.StatusForbidden, "error.html", gin.H{
			"code":    http.StatusForbidden,
			"message": "User is not logged in.",
		})
	}

	c.HTML(http.StatusOK, "user-profile.html", gin.H{
		"username": val,
	})
}
