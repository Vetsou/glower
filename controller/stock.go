package controller

import (
	"database/sql"
	"glower/controller/internal"
	"glower/database/model"
	"glower/database/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGetFlowers(factory repository.StockRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := factory(c)

		flowers, err := repo.GetFlowers()
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Failed to load products. Please try again later.")
			return
		}

		c.HTML(http.StatusOK, "shop-stock.html", gin.H{
			"flowers": flowers,
		})
	}
}

func CreateAddFlower(factory repository.StockRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Name          string   `form:"name" binding:"required"`
			Price         float64  `form:"price" binding:"required"`
			Available     bool     `form:"available"`
			Description   string   `form:"description"`
			DiscountPrice *float64 `form:"discount"`
			Stock         uint     `form:"stock" binding:"required"`
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
			DiscountPrice: sql.NullFloat64{Float64: float64(0), Valid: false},
		}

		if request.DiscountPrice != nil {
			flower.DiscountPrice = sql.NullFloat64{
				Float64: *request.DiscountPrice,
				Valid:   true,
			}
		}

		repo := factory(c)
		flower, err := repo.AddAndGetFlower(&flower, request.Stock)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Failed to add new flower.")
			return
		}

		c.HTML(http.StatusOK, "flower-item.html", flower)
	}
}

func CreateRemoveFlower(factory repository.StockRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		flowerId, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Wrong flower ID.")
			return
		}

		repo := factory(c)

		if err := repo.RemoveFlower(uint(flowerId)); err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Error deleting flower. Please try again later.")
			return
		}

		c.Status(http.StatusOK)
	}
}
