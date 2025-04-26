package repository

import (
	"glower/database/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthRepoFactory func(c *gin.Context) AuthRepository

func CreateAuthRepoFactory() AuthRepoFactory {
	return func(c *gin.Context) AuthRepository {
		tx := c.MustGet("tx").(*gorm.DB)
		return newAuthRepo(tx)
	}
}

type AuthRepository interface {
	InsertUser(user model.User) error
	GetUser(email string) (model.User, error)
}

type authRepo struct {
	db *gorm.DB
}

func newAuthRepo(tx *gorm.DB) AuthRepository {
	return &authRepo{db: tx}
}

func (r *authRepo) InsertUser(user model.User) error {
	return r.db.Create(&user).Error
}

func (r *authRepo) GetUser(email string) (model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
