package controller

import (
	"errors"
	"glower/auth"
	"glower/model"
	"net/http"

	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

type RegisterUserFrom struct {
	FirstName       string `form:"first-name" binding:"required"`
	LastName        string `form:"last-name" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
}

func (r *RegisterUserFrom) validate() error {
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

func RegisterUser(c *gin.Context) {
	var formData RegisterUserFrom
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	if err := formData.validate(); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	if err := passwordvalidator.Validate(formData.Password, minEntropyBits); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Password is too weak. Please use at least 10 characters, including uppercase, lowercase, numbers, and special characters.",
		})
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(formData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "An unexpected error occurred while processing your request.",
		})
		return
	}

	user := model.User{
		FirstName:    formData.FirstName,
		LastName:     formData.LastName,
		Email:        formData.Email,
		PasswordHash: hashedPass,
	}

	if err := model.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Error inserting user to database.",
		})
		return
	}

	c.HTML(http.StatusOK, "register-success.html", gin.H{
		"name": user.FirstName + " " + user.LastName,
	})
}

func Login(c *gin.Context) {
	var request struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&request); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	var user model.User
	if err := model.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid email or password.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid email or password.",
		})
		return
	}

	accessToken, err := auth.CreateJWT(user, user.Email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate token.",
		})
		return
	}

	refreshToken, err := auth.CreateRefreshToken(user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate token.",
		})
		return
	}

	auth.SetCookies(c, &refreshToken, &accessToken)

	c.HTML(http.StatusOK, "login-success.html", gin.H{
		"name": user.FirstName + " " + user.LastName,
	})
}

func Logout(c *gin.Context) {
	auth.CleanCookies(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": "You have successfully logged out.",
	})
}
