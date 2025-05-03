package mocks

import (
	"database/sql"
	"fmt"
	"glower/database/model"
	"net/url"
	"strconv"

	"github.com/stretchr/testify/mock"
)

type StockRepoMock struct{ mock.Mock }

func (m *StockRepoMock) GetFlowers() ([]model.Flower, error) {
	args := m.Called()
	return args.Get(0).([]model.Flower), args.Error(1)
}

func (m *StockRepoMock) AddFlower(flower model.Flower, flowerStock uint) error {
	args := m.Called(flower, flowerStock)
	return args.Error(0)
}

func (m *StockRepoMock) RemoveFlower(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func GetFlowers() []model.Flower {
	return []model.Flower{
		{
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

func GetOptionalFieldsAddFlowerForm() url.Values {
	form := url.Values{}
	form.Add("name", "FlowerName")
	form.Add("price", fmt.Sprintf("%f", 12.00))
	form.Add("stock", strconv.FormatUint(uint64(13), 10))

	return form
}
