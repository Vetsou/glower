package controller

import (
	"glower/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetFlowers(c *gin.Context) {
	var flowers []model.Flower

	err := model.DB.Model(&model.Flower{}).Preload("Inventory").Find(&flowers).Error

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to load flowers. Please try again later.",
		})
		return
	}

	c.HTML(http.StatusOK, "shop-stock.html", gin.H{
		"flowers": flowers,
	})
}

func AddFlower(c *gin.Context) {
	var request struct {
		Name          string  `form:"name" binding:"required"`
		Price         float32 `form:"price" binding:"required"`
		Available     bool    `form:"available"`
		Description   string  `form:"description"`
		DiscountPrice float32 `form:"discount"`
		Stock         uint    `form:"stock" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	flower := model.Flower{
		Name:          request.Name,
		Price:         request.Price,
		Available:     request.Available,
		Description:   request.Description,
		DiscountPrice: request.DiscountPrice,
	}

	tx := model.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Database error.",
			})
		}
	}()

	if err := tx.Create(&flower).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add flower to the database.",
		})
		return
	}

	inventory := model.Inventory{
		FlowerID: flower.ID,
		Stock:    request.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add inventory for the flower.",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to commit changes.",
		})
		return
	}

	c.HTML(http.StatusOK, "stock-add.html", gin.H{
		"ID":            flower.ID,
		"name":          flower.Name,
		"price":         flower.Price,
		"available":     flower.Available,
		"description":   flower.Description,
		"discountPrice": flower.DiscountPrice,
		"stock":         inventory.Stock,
	})
}

func RemoveFlower(c *gin.Context) {
	id := c.Param("id")

	if err := model.DB.Select(clause.Associations).Delete(&model.Flower{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error deleting flower from DB.",
		})
		return
	}

	c.Status(200)
}
