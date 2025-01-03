package model

import (
	"errors"
)

type AddFlowerForm struct {
	Name          string  `form:"name" binding:"required"`
	Price         float32 `form:"price" binding:"required"`
	Available     bool    `form:"available"`
	Description   string  `form:"description"`
	DiscountPrice float32 `form:"discount"`
	Stock         uint    `form:"stock" binding:"required"`
}

type RegisterUserFrom struct {
	FirstName       string `form:"first-name" binding:"required"`
	LastName        string `form:"last-name" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
}

func (r *RegisterUserFrom) Validate() error {
	if len(r.FirstName) > 50 {
		return errors.New("first name cannot be longer than 50 characters")
	}

	if len(r.LastName) > 50 {
		return errors.New("last name cannot be longer than 50 characters")
	}

	if len(r.Email) > 70 {
		return errors.New("email cannot be longer than 70 characters")
	}

	if len(r.Password) > 60 {
		return errors.New("password cannot be longer than 60 characters")
	}

	if len(r.ConfirmPassword) > 60 {
		return errors.New("confirm password cannot be longer than 60 characters")
	}

	if r.Password != r.ConfirmPassword {
		return errors.New("passwords must match")
	}

	return nil
}
