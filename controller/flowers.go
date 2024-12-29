package controller

import (
	"glower/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFlowers(c *gin.Context) {
	var flowers []model.Flower

	err := model.DB.Select("ID", "Name", "Price", "Available", "Description", "Discount").Find(&flowers).Error
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
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	availableStr := c.PostForm("available")
	description := c.PostForm("description")
	discountStr := c.PostForm("discount")

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid price format.",
		})
		return
	}

	var available bool
	if availableStr == "" {
		available = false
	} else if availableStr == "true" {
		available = true
	} else {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid available format. Should be 'true' or empty for 'false'.",
		})
		return
	}

	discount, err := strconv.ParseFloat(discountStr, 32)
	if err != nil {
		discount = 0
	}

	flower := model.Flower{
		Name:        name,
		Price:       float32(price),
		Available:   available,
		Description: description,
		Discount:    float32(discount),
	}

	err = model.DB.Create(&flower).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to add flower to the database.",
		})
		return
	}

	c.HTML(http.StatusOK, "stock-add.html", gin.H{
		"ID":          flower.ID,
		"name":        flower.Name,
		"price":       flower.Price,
		"available":   flower.Available,
		"description": flower.Description,
		"discount":    flower.Discount,
	})
}

func RemoveFlower(c *gin.Context) {
	id := c.Param("id")

	if err := model.DB.Delete(&model.Flower{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error deleting flower from DB.",
		})
		return
	}

	c.Status(200)
}
