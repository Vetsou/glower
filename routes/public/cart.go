package public

import (
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCart(e *gin.Engine, db *gorm.DB) {
	flowersGroup := e.Group("/cart",
		middleware.CreateAuth(true),
		middleware.CreateRolesAuth(model.RoleUser),
		middleware.CreateTransaction(db))

	factory := repository.CreateCartRepoFactory()

	flowersGroup.GET("/", controller.CreateGetCartItems(factory))
	flowersGroup.POST("/", controller.CreateAddCartItem(factory))
	flowersGroup.DELETE("/:id", controller.CreateRemoveCartItem(factory))
}
