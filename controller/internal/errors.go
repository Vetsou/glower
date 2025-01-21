package internal

import (
	"github.com/gin-gonic/gin"
)

func SetPartialError(c *gin.Context, code int, msg string) {
	c.HTML(code, "error-alert.html", gin.H{"errorMessage": msg})
}
