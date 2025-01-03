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
	var formData model.AddFlowerForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	flower := model.Flower{
		Name:          formData.Name,
		Price:         formData.Price,
		Available:     formData.Available,
		Description:   formData.Description,
		DiscountPrice: formData.DiscountPrice,
	}

	tx := model.DB.Begin()

	if err := tx.Create(&flower).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add flower to the database.",
		})
		return
	}

	inventory := model.Inventory{
		FlowerID: flower.ID,
		Stock:    formData.Stock,
	}

	if err := tx.Create(&inventory).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add inventory for the flower.",
		})
		return
	}

	tx.Commit()

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
