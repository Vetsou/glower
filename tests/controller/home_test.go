package controller

import (
	"glower/controller"
	"glower/initializers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup

func setupHomeRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/")
	group.GET("/", controller.CreateGetHomePage())

	return r
}

// Tests

func TestGetHomePage_NoOper(t *testing.T) {
	// Setup
	router := setupHomeRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotContains(t, resp.Body.String(), "logged out.")
	assert.NotContains(t, resp.Body.String(), "logged in.")
	assert.NotContains(t, resp.Body.String(), "Registration successful.")
}

func TestGetHomePage_Logout(t *testing.T) {
	// Setup
	router := setupHomeRouter()
	req, _ := http.NewRequest("GET", "/?oper=logout", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "You have successfully logged out.")
}

func TestGetHomePage_Login(t *testing.T) {
	// Setup
	router := setupHomeRouter()
	req, _ := http.NewRequest("GET", "/?oper=login", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "You have successfully logged in.")
}

func TestGetHomePage_Register(t *testing.T) {
	// Setup
	router := setupHomeRouter()
	req, _ := http.NewRequest("GET", "/?oper=register", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "Registration successful. Please log in to continue.")
}
