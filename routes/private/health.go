package private

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHealth(e *gin.Engine, db *gorm.DB) {
	healthRoute := e.Group("/health")
	healthRoute.GET("/", controller.CheckHealth(db))
}
