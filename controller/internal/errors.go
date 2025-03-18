package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetPartialError(c *gin.Context, code int, msg string) {
	c.HTML(code, "error-alert.html", gin.H{"errorMessage": msg})
}

func HandlePanic(c *gin.Context, tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		SetPartialError(c, http.StatusInternalServerError, "Internal server error.")
	}
}
