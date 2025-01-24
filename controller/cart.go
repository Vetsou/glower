package controller

import (
	"errors"
	"fmt"
	"glower/controller/internal"
	"glower/database"
	"glower/database/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getUserCart(c *gin.Context, tx *gorm.DB) (model.Cart, error) {
	userId := c.GetUint("id")
	if userId == 0 {
		return model.Cart{}, errors.New("incorrect user id")
	}

	var cart model.Cart
	if err := tx.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = model.Cart{
				UserID: userId,
				Items:  []model.CartItem{},
			}

			if err := tx.Create(&cart).Error; err != nil {
				return model.Cart{}, fmt.Errorf("failed to create cart: %w", err)
			}
		} else {
			return model.Cart{}, fmt.Errorf("failed to fetch cart: %w", err)
		}
	}

	return cart, nil
}

func addOrUpdateCartItem(tx *gorm.DB, cartID uint, flowerID uint) (model.CartItem, error) {
	var cartItem model.CartItem
	err := tx.Model(&model.CartItem{}).
		Where("cart_id = ? AND flower_id = ?", cartID, flowerID).
		First(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = model.CartItem{
				CartID:   cartID,
				FlowerID: flowerID,
				Quantity: 1,
			}

			if err := tx.Create(&cartItem).Error; err != nil {
				return model.CartItem{}, fmt.Errorf("failed to create cart item: %w", err)
			}
		} else {
			return model.CartItem{}, fmt.Errorf("failed to query cart item: %w", err)
		}
	} else {
		cartItem.Quantity++
		if err := tx.Save(&cartItem).Error; err != nil {
			return model.CartItem{}, fmt.Errorf("failed to update cart item: %w", err)
		}
	}

	return cartItem, nil
}

func GetCartItems(c *gin.Context) {
	tx := database.Handle.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			internal.SetPartialError(c, http.StatusInternalServerError, "Internal server error.")
		}
	}()

	cart, err := getUserCart(c, tx)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	var cartItems []model.CartItem
	if err := tx.Preload("Flower").Where("cart_id = ?", cart.ID).Find(&cartItems).Error; err != nil {
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

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			internal.SetPartialError(c, http.StatusInternalServerError, "Internal server error.")
		}
	}()

	var flower model.Flower
	err := tx.Model(&model.Flower{}).Preload("Inventory").Find(&flower, request.FlowerID).Error
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

	cart, err := getUserCart(c, tx)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	cartItem, err := addOrUpdateCartItem(tx, cart.ID, flower.ID)
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

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			internal.SetPartialError(c, http.StatusInternalServerError, "Internal server error.")
		}
	}()

	cart, err := getUserCart(c, tx)
	if err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusInternalServerError, "Failed to retrieve user cart.")
		return
	}

	var cartItem model.CartItem
	if err := tx.Where("id = ? AND cart_id = ?", cartItemId, cart.ID).First(&cartItem).Error; err != nil {
		tx.Rollback()
		internal.SetPartialError(c, http.StatusNotFound, "Cart item not found.")
		return
	}

	if err := tx.Delete(&cartItem).Error; err != nil {
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
