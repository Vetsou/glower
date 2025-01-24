package auth

import (
	"errors"
	"fmt"
	"glower/database/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RefreshTokenName = "refresh-token"
	AccessTokenName  = "access-token"
	DomainName       = "localhost"
)

var (
	jwtSecret     = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))
)

type UserTokenData struct {
	Id   uint
	User string
}

func CreateJWT(user model.User, email string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"data": map[string]string{
			"id":   fmt.Sprintf("%d", user.ID),
			"user": user.FirstName + " " + user.LastName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func CreateRefreshToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user": user.FirstName + " " + user.LastName,
		"exp":  time.Now().Add(3 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func VerifyToken(tokenStr string) (*jwt.MapClaims, error) {
	decodedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !decodedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := decodedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}

	return &claims, nil
}

func GetUserClaims(claims *jwt.MapClaims) (UserTokenData, error) {
	tokenData := (*claims)["data"].(map[string]interface{})

	userId, err := strconv.ParseUint(tokenData["id"].(string), 10, 32)
	if err != nil {
		return UserTokenData{}, err
	}

	userData := UserTokenData{
		Id:   uint(userId),
		User: tokenData["user"].(string),
	}

	return userData, nil
}
