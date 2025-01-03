package controller

import (
	"glower/model"
	"net/http"

	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

func RegisterUser(c *gin.Context) {
	var formData model.RegisterUserFrom
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	if err := formData.Validate(); err != nil {
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
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Password hashing failed. Please try again later.",
			})
		default:
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"code":    http.StatusInternalServerError,
				"message": "An unexpected error occurred while processing your request.",
			})
		}
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

}
