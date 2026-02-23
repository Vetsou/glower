package mocks

import (
	"database/sql"
	"fmt"
	"glower/database/model"
	"net/url"
	"strconv"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type StockRepoMock struct{ mock.Mock }

func (m *StockRepoMock) GetFlowers() ([]model.Flower, error) {
	args := m.Called()
	return args.Get(0).([]model.Flower), args.Error(1)
}

func (m *StockRepoMock) AddAndGetFlower(flower *model.Flower, flowerStock uint) (model.Flower, error) {
	args := m.Called(flower, flowerStock)
	return args.Get(0).(model.Flower), args.Error(1)
}

func (m *StockRepoMock) RemoveFlower(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func GetFlowers() []model.Flower {
	return []model.Flower{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "Sunflower",
			Price:       9.99,
			Available:   false,
			Description: "Yellow flower",
			Inventory: model.Inventory{
				FlowerID: 1,
				Stock:    10,
			},
		},
		{
			Model:         gorm.Model{ID: 2},
			Name:          "Poppy",
			Price:         7.99,
			Available:     true,
			Description:   "Red flower",
			DiscountPrice: sql.NullFloat64{Float64: 5.99, Valid: true},
			Inventory: model.Inventory{
				FlowerID: 1,
				Stock:    13,
			},
		},
	}
}

func GetValidAddFlowerForm() url.Values {
	form := url.Values{}
	form.Add("name", "FlowerName")
	form.Add("price", fmt.Sprintf("%f", 15.00))
	form.Add("available", strconv.FormatBool(true))
	form.Add("description", "Nice flower")
	form.Add("discount", fmt.Sprintf("%f", 10.0))
	form.Add("stock", strconv.FormatUint(uint64(12), 10))

	return form
}

func GetValidAddFlowerModel() model.Flower {
	return model.Flower{
		Model:         gorm.Model{ID: 67},
		Name:          "FlowerName",
		Price:         15.00,
		Available:     true,
		Description:   "Nice flower",
		DiscountPrice: sql.NullFloat64{Float64: 10.0, Valid: true},
		Inventory: model.Inventory{
			FlowerID: 1,
			Stock:    12,
		},
	}
}

func GetOptionalFieldsAddFlowerForm() url.Values {
	form := url.Values{}
	form.Add("name", "FlowerName")
	form.Add("price", fmt.Sprintf("%f", 12.00))
	form.Add("stock", strconv.FormatUint(uint64(13), 10))

	return form
}

func GetOptionalFieldsAddFlowerModel() model.Flower {
	return model.Flower{
		Model: gorm.Model{ID: 43},
		Name:  "FlowerName",
		Price: 12.00,
		Inventory: model.Inventory{
			FlowerID: 1,
			Stock:    13,
		},
	}
}
