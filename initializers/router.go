package initializers

import (
	"glower/routes/private"
	"glower/routes/public"

	"github.com/gin-gonic/gin"
)

func RegisterServiceRoutes(e *gin.Engine) {
	public.RegisterHome(e)
	public.RegisterFlowers(e)
	public.RegisterUser(e)
	public.RegisterAuth(e)
}

func RegisterPrivateRoutes(e *gin.Engine) {
	reg := CreateMetricsRegistry()
	private.RegisterMetrics(e, reg)
}
