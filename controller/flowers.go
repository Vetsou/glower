package controller

import (
	"glower/controller/internal"
	"glower/database"
	"glower/database/model"
	"glower/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFlowers(c *gin.Context) {
	repo := repository.NewStockRepo(database.Handle)

	flowers, err := repo.GetFlowers()
	if err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to load products. Please try again later.")
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
		internal.SetPartialError(c, http.StatusBadRequest, "Invalid form data. Please fill all required fields.")
		return
	}

	flower := model.Flower{
		Name:          request.Name,
		Price:         request.Price,
		Available:     request.Available,
		Description:   request.Description,
		DiscountPrice: request.DiscountPrice,
	}

	tx := database.Handle.Begin()
	defer internal.HandlePanic(c, tx)
	repo := repository.NewStockRepo(tx)

	if err := repo.AddFlower(flower, request.Stock); err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to add new flower.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "We couldn't save your changes. Please try again.")
		return
	}

	c.HTML(http.StatusOK, "stock-add.html", gin.H{
		"ID":            flower.ID,
		"name":          flower.Name,
		"price":         flower.Price,
		"available":     flower.Available,
		"description":   flower.Description,
		"discountPrice": flower.DiscountPrice,
		"stock":         request.Stock,
	})
}

func RemoveFlower(c *gin.Context) {
	repo := repository.NewStockRepo(database.Handle)

	if err := repo.RemoveFlower(c.Param("id")); err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "Error deleting flower. Please try again later.")
		return
	}

	c.Status(http.StatusOK)
}
