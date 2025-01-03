package model

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	FlowerID uint `gorm:"not null"`
	Stock    uint `gorm:"not null"`
}
