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
