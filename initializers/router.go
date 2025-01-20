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
	public.RegisterCart(e)
}

func RegisterPrivateRoutes(e *gin.Engine) {
	private.RegisterMetrics(e)
}
