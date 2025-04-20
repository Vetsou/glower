package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomePage(c *gin.Context) {
	oper := c.Query("oper")
	var message string

	switch oper {
	case "logout":
		message = "You have successfully logged out."
	case "login":
		message = "You have successfully logged in."
	case "register":
		message = "Registration successful. Please log in to continue."
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": message,
	})
}
