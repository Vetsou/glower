package main

import (
	"glower/initializers"
	"glower/model"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	model.InitDatabase()
}

func main() {
	e := gin.Default()
	initializers.LoadHTMLTemplates(e)
	initializers.RegisterServerRoutes(e)

	e.Run()
}
