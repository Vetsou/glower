package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); exists {
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		c.HTML(http.StatusOK, "user-register.html", nil)
	}
}

func CreateLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); exists {
			c.Redirect(http.StatusFound, "/user/profile")
			return
		}

		c.HTML(http.StatusOK, "user-login.html", nil)
	}
}

func CreateProfilePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user")

		if !exists {
			c.HTML(http.StatusUnauthorized, "error-page.html", gin.H{
				"code":    http.StatusUnauthorized,
				"message": "User is not logged in.",
			})
			return
		}

		c.HTML(http.StatusOK, "user-profile.html", gin.H{
			"username": val,
		})
	}
}
