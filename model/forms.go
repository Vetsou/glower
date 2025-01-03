package model

type AddFlowerForm struct {
	Name          string  `form:"name" binding:"required"`
	Price         float32 `form:"price" binding:"required"`
	Available     bool    `form:"available"`
	Description   string  `form:"description"`
	DiscountPrice float32 `form:"discount"`
	Stock         uint    `form:"stock" binding:"required"`
}
