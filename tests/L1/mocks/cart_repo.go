package mocks

import (
	"database/sql"
	"glower/database/model"

	"github.com/stretchr/testify/mock"
)

type CartRepoMock struct{ mock.Mock }

func (r *CartRepoMock) GetUserCart(userId uint) (model.Cart, error) {
	args := r.Called(userId)
	return args.Get(0).(model.Cart), args.Error(1)
}

func (r *CartRepoMock) AddOrUpdateCartItem(cartID, flowerID uint) (model.CartItem, error) {
	args := r.Called(cartID, flowerID)
	return args.Get(0).(model.CartItem), args.Error(1)
}

func (r *CartRepoMock) GetCartItems(cartID uint) ([]model.CartItem, error) {
	args := r.Called(cartID)
	return args.Get(0).([]model.CartItem), args.Error(1)
}

func (r *CartRepoMock) RemoveCartItem(cartID uint, cartItemID uint) error {
	args := r.Called(cartID, cartItemID)
	return args.Error(0)
}

func (r *CartRepoMock) DecreaseInventoryAndGetFlower(flowerID uint) (model.Flower, error) {
	args := r.Called(flowerID)
	return args.Get(0).(model.Flower), args.Error(1)
}

func GetEmptyCartItems() []model.CartItem {
	return []model.CartItem{}
}

func GetCartFlower(isAvailable bool) model.Flower {
	return model.Flower{
		Name:        "Sunflower",
		Price:       12.99000,
		Available:   isAvailable,
		Description: "Nice flower",
		Inventory: model.Inventory{
			FlowerID: 13,
			Stock:    10,
		},
	}
}

func GetAddedCartFlower() model.CartItem {
	return model.CartItem{
		CartID:   1,
		FlowerID: 13,
		Flower: model.Flower{
			Name:        "Sunflower",
			Price:       12.99000,
			Available:   true,
			Description: "Nice flower",
			Inventory: model.Inventory{
				FlowerID: 13,
				Stock:    10,
			},
		},
		Quantity: 1,
	}
}

func GetTestCartItems() []model.CartItem {
	return []model.CartItem{
		{
			CartID:   1,
			FlowerID: 11,
			Flower: model.Flower{
				Name:        "Sunflower",
				Price:       9.99000,
				Available:   false,
				Description: "Yellow flower",
				Inventory: model.Inventory{
					FlowerID: 1,
					Stock:    10,
				},
			},
			Quantity: 1,
		},
		{
			CartID:   2,
			FlowerID: 12,
			Flower: model.Flower{
				Name:          "Poppy",
				Price:         7.99000,
				Available:     true,
				Description:   "Red flower",
				DiscountPrice: sql.NullFloat64{Float64: 5.99000, Valid: true},
				Inventory: model.Inventory{
					FlowerID: 1,
					Stock:    13,
				},
			},
			Quantity: 2,
		},
	}
}
