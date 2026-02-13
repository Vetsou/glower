package public

import (
	"glower/controller"

	"github.com/gin-gonic/gin"
)

func RegisterStaticFiles(e *gin.Engine) {
	e.GET("/favicon.ico", controller.CreateServeFavicon())
}
