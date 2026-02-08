package model

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	gorm.Model
	FirstName    string `gorm:"size:50;not null"`
	LastName     string `gorm:"size:50;not null"`
	Email        string `gorm:"size:70;not null;unique"`
	PasswordHash []byte `gorm:"size:60;not null"`
	Role         Role   `gorm:"not null"`
}
