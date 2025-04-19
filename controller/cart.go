package controller

import (
	"fmt"
	"glower/controller/internal"
	"glower/database"
	"glower/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCartItems(c *gin.Context) {
	tx := database.Handle.Begin()
	defer internal.HandlePanic(c, tx)
	repo := repository.NewCartRepo(tx)

	cart, err := repo.GetUserCart(c.GetUint("id"))
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	cartItems, err := repo.GetCartItems(cart.ID)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve cart items.")
		return
	}

	var totalPrice float32
	for _, item := range cartItems {
		totalPrice += float32(item.Quantity) * item.Flower.Price
	}

	if err := tx.Commit().Error; err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to commit transaction.")
		return
	}

	c.HTML(http.StatusOK, "user-cart.html", gin.H{
		"cartItems":  cartItems,
		"totalPrice": totalPrice,
	})
}

func AddCartItem(c *gin.Context) {
	var request struct {
		FlowerID uint `form:"flowerId"`
	}

	if err := c.ShouldBind(&request); err != nil {
		internal.SetPartialError(c, http.StatusBadRequest, "Invalid request body.")
		return
	}

	tx := database.Handle.Begin()
	defer internal.HandlePanic(c, tx)
	repo := repository.NewCartRepo(tx)

	flower, err := repo.GetFlowerByID(request.FlowerID)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusNotFound, "Flower not found.")
		return
	}

	if !flower.Available {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusBadRequest, "This flower in not available anymore.")
		return
	}

	cart, err := repo.GetUserCart(c.GetUint("id"))
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	cartItem, err := repo.AddOrUpdateCartItem(cart.ID, flower.ID)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to add flower to cart.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to commit transaction.")
		return
	}

	c.HTML(http.StatusOK, "success-alert.html", gin.H{
		"message": fmt.Sprintf(
			"Flower %s was added to your cart. You currently have %d %s in your cart.",
			flower.Name, cartItem.Quantity, flower.Name),
	})
}

func RemoveCartItem(c *gin.Context) {
	cartItemId := c.Param("id")

	if cartItemId == "" {
		internal.SetPartialError(c, http.StatusBadRequest, "Cart item ID is required.")
		return
	}

	tx := database.Handle.Begin()
	defer internal.HandlePanic(c, tx)
	repo := repository.NewCartRepo(tx)

	cart, err := repo.GetUserCart(c.GetUint("id"))
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	if err := repo.RemoveCartItem(cart.ID, cartItemId); err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to remove cart item.")
		return
	}

	if err := tx.Commit().Error; err != nil {
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to commit transaction.")
		return
	}

	c.HTML(http.StatusOK, "success-alert.html", gin.H{
		"message": "Item was removed from your cart.",
	})
}
