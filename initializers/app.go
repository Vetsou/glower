package initializers

import (
	"glower/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp() *App {
	db := database.Init()

	router := gin.Default()

	InitHTMLTemplates(router, "")
	RegisterServiceMiddleware(router)
	RegisterServiceRoutes(router, db)

	return &App{
		Router: router,
		DB:     db,
	}
}
