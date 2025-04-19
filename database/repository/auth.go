package repository

import (
	"glower/database/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	InsertUser(user model.User) error
	GetUser(email string) (model.User, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(tx *gorm.DB) AuthRepository {
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
