package main

import (
	"glower/database"
	"glower/initializers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	initializers.LoadEnvVariables()
	db = database.Init()
}

func main() {
	gin.SetMode(gin.DebugMode)

	publicRouter := gin.Default()
	initializers.RegisterServiceMiddleware(publicRouter)
	initializers.InitHTMLTemplates(publicRouter, "")
	initializers.RegisterServiceRoutes(publicRouter, db)

	// Run private router
	go func() {
		privateRouter := gin.New()
		initializers.RegisterPrivateRoutes(privateRouter)

		if err := privateRouter.Run(os.Getenv("PRIV_ADDR")); err != nil {
			log.Fatalf("Failed to start private metrics server: %v", err)
		}
	}()

	// Run public router
	if err := publicRouter.Run(); err != nil {
		log.Fatalf("Failed to start public server: %v", err)
	}
}
