package model

import "gorm.io/gorm"

type Flower struct {
	gorm.Model
	Name        string
	Price       float32
	Available   bool
	Description string
	Discount    float32
}

type Stock struct {
	gorm.Model
	ProductID uint
	Count     uint
	Limit     uint
}
