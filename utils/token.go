package utils

import (
	"errors"
	"glower/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

func CreateJWT(user model.User, email string) (string, error) {
	claims := jwt.MapClaims{
		"user": user.FirstName + " " + user.LastName,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
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
