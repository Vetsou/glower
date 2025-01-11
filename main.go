package main

import (
	"glower/initializers"
	"glower/model"
	"glower/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	model.InitDatabase()
}

func main() {
	e := gin.Default()
	initializers.InitHTMLTemplates(e)
	routes.RegisterServiceRoutes(e)

	e.Run()
}
