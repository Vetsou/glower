package routes

import "github.com/gin-gonic/gin"

func RegisterServiceRoutes(e *gin.Engine) {
	registerHome(e)
	registerFlowers(e)
	registerUser(e)
	registerAuth(e)
}
