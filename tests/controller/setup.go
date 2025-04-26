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
	}

	return auth.CreateJWT(user, "test@example.com")
}
