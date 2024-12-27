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
	e.LoadHTMLGlob("templates/*")

	initializers.RegisterServerRoutes(e)

	e.Run()
}
