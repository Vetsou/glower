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
	database.Init()
}

func main() {
	publicRouter := gin.Default()
	initializers.RegisterServiceMiddleware(publicRouter)
	initializers.InitHTMLTemplates(publicRouter)
	initializers.RegisterServiceRoutes(publicRouter)

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
