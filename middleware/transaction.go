package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handlePanic(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}

func isCorrectResponse(status int) bool {
	if status == http.StatusOK {
		return true
	}

	if status == http.StatusCreated {
		return true
	}

	return false
}

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()
		defer handlePanic(tx)

		c.Set("tx", tx)
		c.Next()

		if isCorrectResponse(c.Writer.Status()) {
			if err := tx.Commit().Error; err != nil {
				c.HTML(http.StatusInternalServerError, "error-alert.html", gin.H{
					"errorMessage": "We couldn't save your changes. Please try again.",
				})
			}
		} else {
			tx.Rollback()
		}
	}
}
