package repository

import (
	"errors"
	"fmt"
	"glower/database/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetUserCart(userId uint) (model.Cart, error)
	AddOrUpdateCartItem(cartID, flowerID uint) (model.CartItem, error)
	GetCartItems(cartID uint) ([]model.CartItem, error)
	RemoveCartItem(cartID uint, cartItemID string) error
	GetFlowerByID(flowerID uint) (model.Flower, error)
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepository(tx *gorm.DB) CartRepository {
	return &cartRepo{db: tx}
}

func (r *cartRepo) GetUserCart(userId uint) (model.Cart, error) {
	if userId == 0 {
		return model.Cart{}, errors.New("incorrect user id")
	}

	var cart model.Cart
	if err := r.db.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = model.Cart{
				UserID: userId,
				Items:  []model.CartItem{},
			}

			if err := r.db.Create(&cart).Error; err != nil {
				return model.Cart{}, fmt.Errorf("failed to create cart: %w", err)
			}
		} else {
			return model.Cart{}, fmt.Errorf("failed to fetch cart: %w", err)
		}
	}

	return cart, nil
}

func (r *cartRepo) AddOrUpdateCartItem(cartID, flowerID uint) (model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Model(&model.CartItem{}).
		Where("cart_id = ? AND flower_id = ?", cartID, flowerID).
		First(&cartItem).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cartItem = model.CartItem{
				CartID:   cartID,
				FlowerID: flowerID,
				Quantity: 1,
			}

			if err := r.db.Create(&cartItem).Error; err != nil {
				return model.CartItem{}, fmt.Errorf("failed to create cart item: %w", err)
			}
		} else {
			return model.CartItem{}, fmt.Errorf("failed to query cart item: %w", err)
		}
	} else {
		cartItem.Quantity++
		if err := r.db.Save(&cartItem).Error; err != nil {
			return model.CartItem{}, fmt.Errorf("failed to update cart item: %w", err)
		}
	}

	return cartItem, nil
}

func (r *cartRepo) GetCartItems(cartID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	if err := r.db.Preload("Flower").Where("cart_id = ?", cartID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *cartRepo) RemoveCartItem(cartID uint, cartItemID string) error {
	return r.db.Where("id = ? AND cart_id = ?", cartItemID, cartID).Delete(&model.CartItem{}).Error
}

func (r *cartRepo) GetFlowerByID(flowerID uint) (model.Flower, error) {
	var flower model.Flower
	err := r.db.Preload("Inventory").First(&flower, flowerID).Error
	return flower, err
}
