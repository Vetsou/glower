package repository

import (
	"glower/database/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockRepoFactory func(c *gin.Context) StockRepository

func CreateStockRepoFactory() StockRepoFactory {
	return func(c *gin.Context) StockRepository {
		tx := c.MustGet("tx").(*gorm.DB)
		return newStockRepo(tx)
	}
}

type StockRepository interface {
	GetFlowers() ([]model.Flower, error)
	AddAndGetFlower(flower *model.Flower, flowerStock uint) (model.Flower, error)
	RemoveFlower(id uint) error
}

type stockRepo struct {
	db *gorm.DB
}

func newStockRepo(tx *gorm.DB) StockRepository {
	return &stockRepo{db: tx}
}

func (r *stockRepo) GetFlowers() ([]model.Flower, error) {
	var flowers []model.Flower

	err := r.db.Model(&model.Flower{}).Preload("Inventory").Find(&flowers).Error
	if err != nil {
		return nil, err
	}

	return flowers, nil
}

func (r *stockRepo) AddAndGetFlower(flower *model.Flower, count uint) (model.Flower, error) {
	if err := r.db.Create(flower).Error; err != nil {
		return model.Flower{}, err
	}

	inventory := model.Inventory{
		FlowerID: flower.ID,
		Stock:    count,
	}

	if err := r.db.Create(&inventory).Error; err != nil {
		return model.Flower{}, err
	}

	var result model.Flower

	err := r.db.
		Preload("Inventory").
		First(&result, flower.ID).
		Error

	if err != nil {
		return model.Flower{}, err
	}

	return result, nil
}

func (r *stockRepo) RemoveFlower(id uint) error {
	err := r.db.
		Select(clause.Associations).
		Delete(&model.Flower{}, id).Error

	if err != nil {
		return err
	}

	return nil
}
