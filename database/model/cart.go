package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint       `gorm:"not null"`
	User   User       `gorm:"not null"`
	Items  []CartItem `gorm:"constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	gorm.Model
	CartID   uint   `gorm:"not null"`
	FlowerID uint   `gorm:"not null"`
	Flower   Flower `gorm:"foreignKey:FlowerID"`
	Quantity uint   `gorm:"not null;default:0"`
}
