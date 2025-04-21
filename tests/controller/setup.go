package controller

import (
	"glower/auth"
	"glower/database/model"
	"glower/initializers"
	"glower/routes/public"
	"os"

	"github.com/gin-gonic/gin"
)

func setupHomeRouter() *gin.Engine {
	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")
	public.RegisterHome(r)
	return r
}

func setupUserRouter() *gin.Engine {
	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")
	public.RegisterUser(r)

	return r
}

func setupTokenAuth() (string, error) {
	os.Setenv("ACCESS_TOKEN_SECRET", "test-access-token-value")
	os.Setenv("REFRESH_TOKEN_SECRET", "test-refresh-token-value")

	user := model.User{
		FirstName: "Test",
		LastName:  "User",
	}

	return auth.CreateJWT(user, "test@example.com")
}
