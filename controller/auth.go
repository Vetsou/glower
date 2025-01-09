package controller

import (
	"glower/auth"
	"glower/input/form"
	"glower/model"
	"net/http"

	"github.com/gin-gonic/gin"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

func RegisterUser(c *gin.Context) {
	var formData form.RegisterUserFrom
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
	var formData form.LoginUserForm
	if err := c.ShouldBind(&formData); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid form data: " + err.Error(),
		})
		return
	}

	var user model.User
	if err := model.DB.Where("email = ?", formData.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid email or password.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(formData.Password)); err != nil {
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

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.RefreshTokenName,
		Value:    refreshToken,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   3 * 24 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.AccessTokenName,
		Value:    accessToken,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   24 * 60 * 60,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.HTML(http.StatusOK, "login-success.html", gin.H{
		"name": user.FirstName + " " + user.LastName,
	})
}

func Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.RefreshTokenName,
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     auth.AccessTokenName,
		Value:    "",
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.HTML(http.StatusOK, "index.html", gin.H{
		"message": "You have successfully logged out.",
	})
}
