package controller

import (
	"glower/auth"
	"glower/database/model"
	"os"
)

func createTokenMock() (string, error) {
	os.Setenv("ACCESS_TOKEN_SECRET", "test-access-token-value")
	os.Setenv("REFRESH_TOKEN_SECRET", "test-refresh-token-value")

	user := model.User{
		FirstName: "Test",
		LastName:  "User",
		Role:      model.RoleAdmin,
	}

	return auth.CreateJWT(user)
}
