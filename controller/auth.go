package controller

import (
	"errors"
	"glower/auth"
	"glower/controller/internal"
	"glower/database/model"
	"glower/database/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

type registerUserFrom struct {
	FirstName       string `form:"first-name" binding:"required"`
	LastName        string `form:"last-name" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
}

func (r *registerUserFrom) validate() error {
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

func CreateRegister(factory repository.AuthRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData registerUserFrom
		if err := c.ShouldBind(&formData); err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Invalid form data. Please fill all required fields.")
			return
		}

		if err := formData.validate(); err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Invalid form data: "+err.Error())
			return
		}

		if err := passwordvalidator.Validate(formData.Password, minEntropyBits); err != nil {
			internal.SetPartialError(c, http.StatusBadRequest,
				"Password is too weak. Please use at least 10 characters, including uppercase, lowercase, numbers, and special characters.")
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(formData.Password), bcrypt.DefaultCost)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "An unexpected error occurred while processing your request.")
			return
		}

		user := model.User{
			FirstName:    formData.FirstName,
			LastName:     formData.LastName,
			Email:        formData.Email,
			PasswordHash: hashedPass,
		}

		repo := factory(c)
		if err := repo.InsertUser(user); err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "Cannot register user. Please try again later.")
			return
		}

		c.Header("HX-Redirect", "/?oper=register")
		c.Status(http.StatusOK)
	}
}

func CrateLogin(factory repository.AuthRepoFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Email    string `form:"email" binding:"required"`
			Password string `form:"password" binding:"required"`
		}

		if err := c.ShouldBind(&request); err != nil {
			internal.SetPartialError(c, http.StatusBadRequest, "Invalid form data. Please fill all required fields.")
			return
		}

		repo := factory(c)
		user, err := repo.GetUser(request.Email)
		if err != nil {
			internal.SetPartialError(c, http.StatusUnauthorized, "Invalid email or password.")
			return
		}

		if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(request.Password)); err != nil {
			internal.SetPartialError(c, http.StatusUnauthorized, "Invalid email or password.")
			return
		}

		accessToken, err := auth.CreateJWT(user, user.Email)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "An internal error occurred while processing your login. Please try again later.")
			return
		}

		refreshToken, err := auth.CreateRefreshToken(user)
		if err != nil {
			internal.SetPartialError(c, http.StatusInternalServerError, "An internal error occurred while processing your login. Please try again later.")
			return
		}

		auth.SetCookies(c, &refreshToken, &accessToken)
		c.Header("HX-Redirect", "/?oper=login")
		c.Status(http.StatusOK)
	}
}

func CreateLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.CleanCookies(c)

		c.Header("HX-Redirect", "/?oper=logout")
		c.Status(http.StatusOK)
	}
}
