package model

import "gorm.io/gorm"

type Flower struct {
	gorm.Model
	Name          string  `gorm:"not null"`
	Price         float32 `gorm:"not null"`
	Available     bool    `gorm:"default:false"`
	Description   string
	DiscountPrice float32   `gorm:"default:0"`
	Inventory     Inventory `gorm:"foreignKey:FlowerID;constraint:OnDelete:CASCADE;"`
}

type Inventory struct {
	gorm.Model
	FlowerID uint `gorm:"not null"`
	Stock    uint `gorm:"not null;default:0"`
}
