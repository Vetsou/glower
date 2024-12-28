package controller

import (
	"glower/model"
)

func GetFlowersStock() ([]model.Flower, error) {
	var flowers []model.Flower

	err := model.DB.Select("Name", "Price", "Available", "Description", "Discount").Find(&flowers).Error
	if err != nil {
		return nil, err
	}

	return flowers, nil
}
