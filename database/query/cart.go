package query

import (
	"errors"
	"fmt"
	"glower/database/model"

	"gorm.io/gorm"
)

func GetUserCart(userId uint, tx *gorm.DB) (model.Cart, error) {
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

func AddOrUpdateCartItem(cartID uint, flowerID uint, tx *gorm.DB) (model.CartItem, error) {
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
