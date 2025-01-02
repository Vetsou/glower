package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseFormBool(input string) (bool, error) {
	if input == "" {
		return false, nil
	}

	if input == "true" {
		return true, nil
	}

	return false, errors.New("invalid boolean format: should be 'true' or empty for 'false'")
}

type addFlowerForm struct {
	name          string
	price         float32
	available     bool
	description   string
	discountPrice float32
	stock         uint
}

func (f *addFlowerForm) fromContext(c *gin.Context) error {
	f.name = c.PostForm("name")

	priceStr := c.PostForm("price")
	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		return fmt.Errorf("invalid price format")
	}
	f.price = float32(price)

	availableStr := c.PostForm("available")
	available, err := parseFormBool(availableStr)
	if err != nil {
		return fmt.Errorf("invalid available format: %v", err)
	}
	f.available = available

	f.description = c.PostForm("description")

	discountStr := c.PostForm("discount")
	discountPrice, err := strconv.ParseFloat(discountStr, 32)
	if err != nil {
		f.discountPrice = 0
	} else {
		f.discountPrice = float32(discountPrice)
	}

	stockStr := c.PostForm("stock")
	stock, err := strconv.ParseUint(stockStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid stock format")
	}
	f.stock = uint(stock)

	return nil
}
