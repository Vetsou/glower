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
	AddFlower(flower model.Flower, flowerStock uint) error
	RemoveFlower(id string) error
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

func (r *stockRepo) AddFlower(flower model.Flower, count uint) error {
	if err := r.db.Create(&flower).Error; err != nil {
		return err
	}

	inventory := model.Inventory{
		FlowerID: flower.ID,
		Stock:    count,
	}

	if err := r.db.Create(&inventory).Error; err != nil {
		return err
	}

	return nil
}

func (r *stockRepo) RemoveFlower(id string) error {
	err := r.db.
		Select(clause.Associations).
		Delete(&model.Flower{}, id).Error

	if err != nil {
		return err
	}

	return nil
}
