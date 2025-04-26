package public

import (
	"glower/controller"
	"glower/database/repository"
	"glower/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuth(e *gin.Engine, db *gorm.DB) {
	authGroup := e.Group("/auth")

	factory := repository.CreateAuthRepoFactory()

	authGroup.POST("/signup", middleware.CreateTransaction(db), controller.CreateRegister(factory))
	authGroup.POST("/login", middleware.CreateTransaction(db), controller.CreateLogin(factory))
	authGroup.POST("/logout", middleware.CreateAuth(true), controller.CreateLogout())
}
