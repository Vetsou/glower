package controller

import (
	"glower/auth"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"glower/middleware"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Mock

type mockAuthRepo struct{}

func (r *mockAuthRepo) InsertUser(user model.User) error { return nil }

func (r *mockAuthRepo) GetUser(email string) (model.User, error) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.DefaultCost)

	user := model.User{
		FirstName:    "FirstNameTest",
		LastName:     "LastNameTest",
		Email:        "test@example.com",
		PasswordHash: hashedPass,
	}
	return user, nil
}

func GetMockAuthRepoFactory() repository.AuthRepoFactory {
	return func(c *gin.Context) repository.AuthRepository {
		return &mockAuthRepo{}
	}
}

// Setup

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/auth")
	factory := GetMockAuthRepoFactory()
	group.POST("/signup", controller.CreateRegister(factory))
	group.POST("/login", controller.CreateLogin(factory))
	group.POST("/logout", middleware.CreateAuth(true), controller.CreateLogout())

	return r
}

// Test Cases

type loginTestData struct {
	email, password string
	expectedCode    int
	expectedMsg     string
}

const (
	WrongCredentialsMsg = "Invalid email or password."
	WrongFormatMsg      = "Invalid form data. Please fill all required fields."
)

var loginTestCases = []loginTestData{
	{"test@example.com", "Password123", http.StatusOK, ""},
	{"test@example.com", "wrongPass", http.StatusUnauthorized, WrongCredentialsMsg},
	{"", "", http.StatusBadRequest, WrongFormatMsg},
	{"test@example.com", "", http.StatusBadRequest, WrongFormatMsg},
	{"", "Password123", http.StatusBadRequest, WrongFormatMsg},
	{"test_example.com", "Password123", http.StatusBadRequest, WrongFormatMsg},
	{"testexample.com", "Password123", http.StatusBadRequest, WrongFormatMsg},
	{"test@example", "Password123", http.StatusBadRequest, WrongFormatMsg},
}

type registerTestData struct {
	firstName, lastName       string
	email                     string
	password, confirmPassword string
	expectedCode              int
	expectedMsg               string
}

var registerTestCases = []registerTestData{
	{
		"fn", "ls", "test@example.com", "password123", "",
		http.StatusBadRequest, WrongFormatMsg,
	},
	{
		"fn", "ls", "test@example.com", "", "password123",
		http.StatusBadRequest, WrongFormatMsg,
	},
	{
		"fn", "", "test@example.com", "password123", "password123",
		http.StatusBadRequest, WrongFormatMsg,
	},
	{
		"", "ls", "test@example.com", "password123", "password123",
		http.StatusBadRequest, WrongFormatMsg,
	},
	{
		"veryveryveryveryveryveryveryveryveryverylongfirstname",
		"ls", "test@example.com", "password", "password",
		http.StatusBadRequest,
		"Invalid form data: first name cannot be longer than 50 characters",
	},
	{
		"firstname", "veryveryveryveryveryveryveryveryveryverylonglastname",
		"test@example.com", "password", "password",
		http.StatusBadRequest,
		"Invalid form data: last name cannot be longer than 50 characters",
	},
	{
		"fs", "ls",
		"veryveryveryveryveryveryveryveryveryveryveryveryverylongtest@example.com",
		"password", "password",
		http.StatusBadRequest,
		"Invalid form data: email cannot be longer than 70 characters",
	},
	{
		"fs", "ls", "test@example.com",
		"veryveryveryveryveryveryveryveryveryveryveryveryverylongpassword",
		"password",
		http.StatusBadRequest,
		"Invalid form data: password cannot be longer than 60 characters",
	},
	{
		"fs", "ls", "test@example.com",
		"password",
		"veryveryveryveryveryveryveryveryveryveryveryveryverylongpassword",
		http.StatusBadRequest,
		"Invalid form data: confirm password cannot be longer than 60 characters",
	},
	{
		"fs", "ls", "test@example.com",
		"password",
		"password_diff",
		http.StatusBadRequest,
		"Invalid form data: passwords must match",
	},
	{
		"fs", "ls", "test@example.com",
		"password",
		"password",
		http.StatusBadRequest,
		"Password is too weak. Please use at least 10 characters",
	},
	{
		"fs", "ls", "test@example.com",
		"Qazxc12345!@#$%",
		"Qazxc12345!@#$%",
		http.StatusOK,
		"",
	},
}

// Tests

func TestLogin_NoData(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/auth/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func TestLogin_WithData(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	for _, testData := range loginTestCases {
		form := url.Values{}
		form.Add("email", testData.email)
		form.Add("password", testData.password)

		req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  auth.AccessTokenName,
			Value: token,
		})

		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, testData.expectedCode, resp.Code)
		assert.Contains(t, resp.Body.String(), testData.expectedMsg)
	}
}

func TestRegister_NoData(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/auth/signup", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func TestRegister_WithData(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	for _, testData := range registerTestCases {
		form := url.Values{}
		form.Add("first-name", testData.firstName)
		form.Add("last-name", testData.lastName)
		form.Add("email", testData.email)
		form.Add("password", testData.password)
		form.Add("confirm-password", testData.confirmPassword)

		req, _ := http.NewRequest("POST", "/auth/signup", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  auth.AccessTokenName,
			Value: token,
		})
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, testData.expectedCode, resp.Code)
		assert.Contains(t, resp.Body.String(), testData.expectedMsg)
	}
}

func TestLogout_WithToken(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "/?oper=logout", resp.Header().Get("HX-Redirect"))
	for _, cookie := range resp.Result().Cookies() {
		if cookie.Name == auth.AccessTokenName {
			assert.True(t, cookie.Value == "" && cookie.MaxAge < 0)
		}
		if cookie.Name == auth.RefreshTokenName {
			assert.True(t, cookie.Value == "" && cookie.MaxAge < 0)
		}
	}
}

func TestLogout_NoToken(t *testing.T) {
	// Arrange
	router := setupAuthRouter()

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "User is not logged in.")
}
