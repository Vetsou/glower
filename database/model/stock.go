package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Flower struct {
	gorm.Model
	Name          string  `gorm:"not null"`
	Price         float64 `gorm:"not null"`
	Available     bool    `gorm:"default:false"`
	Description   string
	DiscountPrice sql.NullFloat64 `gorm:"default:null"`
	Inventory     Inventory       `gorm:"foreignKey:FlowerID;constraint:OnDelete:CASCADE;"`
}

type Inventory struct {
	gorm.Model
	FlowerID uint `gorm:"not null"`
	Stock    uint `gorm:"not null;default:0"`
}
