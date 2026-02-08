package public

import (
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterFlowers(e *gin.Engine, db *gorm.DB) {
	flowersGroup := e.Group("/flowers")

	factory := repository.CreateStockRepoFactory()

	flowersGroup.GET("/",
		middleware.CreateTransaction(db),
		controller.CreateGetFlowers(factory))

	flowersGroup.POST("/",
		middleware.CreateAuth(true),
		middleware.CreateRolesAuth(model.RoleAdmin),
		middleware.CreateTransaction(db),
		controller.CreateAddFlower(factory))

	flowersGroup.DELETE("/:id",
		middleware.CreateAuth(true),
		middleware.CreateRolesAuth(model.RoleAdmin),
		middleware.CreateTransaction(db),
		controller.CreateRemoveFlower(factory))
}
