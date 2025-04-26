package controller

import (
	"fmt"
	"glower/controller/internal"
	"glower/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

		var totalPrice float32
		for _, item := range cartItems {
			totalPrice += float32(item.Quantity) * item.Flower.Price
		}

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

		flower, err := repo.GetFlowerByID(request.FlowerID)
		if err != nil {
			internal.SetPartialError(c, http.StatusNotFound, "The requested flower is unavailable.")
			return
		}

		if !flower.Available {
			internal.SetPartialError(c, http.StatusBadRequest, "This flower is no longer available for purchase.")
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
		cartItemId, exists := c.Params.Get("id")

		if !exists {
			internal.SetPartialError(c, http.StatusBadRequest, "Cart item ID is required.")
			return
		}

		repo := factory(c)

		cart, err := repo.GetUserCart(c.GetUint("id"))
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to load your cart.")
			return
		}

		if err := repo.RemoveCartItem(cart.ID, cartItemId); err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Unable to remove cart item.")
			return
		}

		c.HTML(http.StatusOK, "success-alert.html", gin.H{
			"message": "Item was removed from your cart.",
		})
	}
}
