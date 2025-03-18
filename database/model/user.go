package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"not null;unique"`
	PasswordHash []byte `gorm:"not null"`
}
