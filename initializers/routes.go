package initializers

import (
	"glower/routes"

	"github.com/gin-gonic/gin"
)

func RegisterServerRoutes(e *gin.Engine) {
	routes.RegisterHomeRoute(e)
	routes.RegisterFlowersRoute(e)
}
