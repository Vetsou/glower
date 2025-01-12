package initializers

import (
	"glower/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterServiceMiddleware(e *gin.Engine) {
	e.Use(middleware.CountRequest)
}
