package public

import (
	"glower/controller"
	"glower/database/repository"
	"glower/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCart(e *gin.Engine, db *gorm.DB) {
	flowersGroup := e.Group("/cart")

	factory := repository.CreateCartRepoFactory()

	flowersGroup.GET("/",
		middleware.CreateAuth(true),
		middleware.CreateTransaction(db),
		controller.CreateGetCartItems(factory))

	flowersGroup.POST("/",
		middleware.CreateAuth(true),
		middleware.CreateTransaction(db),
		controller.CreateAddCartItem(factory))

	flowersGroup.DELETE("/:id",
		middleware.CreateAuth(true),
		middleware.CreateTransaction(db),
		controller.CreateRemoveCartItem(factory))
}
