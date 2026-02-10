package main

import (
	"glower/database"
	"glower/initializers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	gin.SetMode(gin.DebugMode)

	db := database.Init()

	publicRouter := gin.Default()
	initializers.InitHTMLTemplates(publicRouter, "")
	initializers.RegisterServiceMiddleware(publicRouter)
	initializers.RegisterServiceRoutes(publicRouter, db)

	// Run private router
	go func() {
		privateRouter := gin.New()
		privateRouter.Use(gin.Recovery())
		initializers.RegisterPrivateRoutes(privateRouter, db)

		if err := privateRouter.Run(os.Getenv("PRIV_ADDR")); err != nil {
			log.Fatalf("Failed to start private metrics server: %v", err)
		}
	}()

	// Run public router
	if err := publicRouter.Run(); err != nil {
		log.Fatalf("Failed to start public server: %v", err)
	}
}
