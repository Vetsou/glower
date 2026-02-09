package initializers

import (
	"glower/routes/private"
	"glower/routes/public"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterServiceRoutes(e *gin.Engine, db *gorm.DB) {
	public.RegisterHome(e)
	public.RegisterUser(e)

	public.RegisterFlowers(e, db)
	public.RegisterAuth(e, db)
	public.RegisterCart(e, db)
}

func RegisterPrivateRoutes(e *gin.Engine, db *gorm.DB) {
	private.RegisterMetrics(e)
	private.RegisterHealth(e, db)
}
