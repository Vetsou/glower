package controller

import (
	"errors"
	"fmt"
	"glower/controller/internal"
	"glower/database/model"
	"glower/database/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func calcCartItemsPrice(cartItems *[]model.CartItem) float64 {
	var totalPrice float64
	for _, item := range *cartItems {
		var unitPrice float64

		if item.Flower.DiscountPrice.Valid {
			unitPrice = item.Flower.DiscountPrice.Float64
		} else {
			unitPrice = item.Flower.Price
		}

		totalPrice += float64(item.Quantity) * unitPrice
	}

	return totalPrice
}

func CreateGetCartItems(factory repository.CartRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := factory(c)

		cart, err := repo.GetUserCart(c.GetUint("id"))
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart.")
			return
		}

		cartItems, err := repo.GetCartItems(cart.ID)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart items.")
			return
		}

		totalPrice := calcCartItemsPrice(&cartItems)

		c.HTML(http.StatusOK, "user-cart.html", gin.H{
			"cartItems":  cartItems,
			"totalPrice": totalPrice,
		})
	}
}

func CreateAddCartItem(factory repository.CartRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			FlowerID uint `form:"flowerId" binding:"required"`
		}

		if err := c.ShouldBind(&request); err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Invalid data provided.")
			return
		}

		repo := factory(c)

		flower, err := repo.DecreaseInventoryAndGetFlower(request.FlowerID)
		if err != nil {
			if errors.Is(err, repository.ErrOutOfStock) {
				internal.SetPartialError(c, http.StatusBadRequest, "The requested flower is out of stock.")
				return
			}

			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load flower data.")
			return
		}

		cart, err := repo.GetUserCart(c.GetUint("id"))
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart.")
			return
		}

		cartItem, err := repo.AddOrUpdateCartItem(cart.ID, flower.ID)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart items.")
			return
		}

		c.HTML(http.StatusOK, "success-alert.html", gin.H{
			"message": fmt.Sprintf(
				"Flower %s was added to your cart. You currently have %d %s in your cart.",
				flower.Name, cartItem.Quantity, flower.Name),
		})
	}
}

func CreateRemoveCartItem(factory repository.CartRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Params.ByName("id")
		cartItemId, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Wrong cart item ID.")
			return
		}

		repo := factory(c)
		cart, err := repo.GetUserCart(c.GetUint("id"))
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart.")
			return
		}

		if err := repo.RemoveCartItem(cart.ID, uint(cartItemId)); err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to remove cart item.")
			return
		}

		c.HTML(http.StatusOK, "success-alert.html", gin.H{
			"message": "Item was removed from your cart.",
		})
	}
}
