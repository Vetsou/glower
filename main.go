package main

import (
	"glower/initializers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	// Init Application
	a := initializers.NewApp()

	// Run private router
	go func() {
		privateRouter := gin.New()
		privateRouter.Use(gin.Recovery())
		initializers.RegisterPrivateRoutes(privateRouter)

		if err := privateRouter.Run(os.Getenv("PRIV_ADDR")); err != nil {
			log.Fatalf("Failed to start private metrics server: %v", err)
		}
	}()

	// Run public router
	if err := a.Router.Run(); err != nil {
		log.Fatalf("Failed to start public server: %v", err)
	}
}
