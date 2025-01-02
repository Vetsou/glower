package model

import (
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	FlowerID uint
	Stock    uint
}
