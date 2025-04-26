package controller

import (
	"glower/auth"
	"glower/controller"
	"glower/initializers"
	"glower/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup

func setupUserRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/user")
	group.GET("/register", middleware.CreateAuth(false), controller.CreateRegisterPage())
	group.GET("/login", middleware.CreateAuth(false), controller.CreateLoginPage())
	group.GET("/profile", middleware.CreateAuth(true), controller.CreateProfilePage())

	return r
}

// Tests

func TestGetRegisterPage_NoUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	req, _ := http.NewRequest("GET", "/user/register", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User Register")
}

func TestGetRegisterPage_WithUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/register", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusFound, resp.Code)
	assert.NotContains(t, resp.Body.String(), "User Register")
}

func TestGetLoginPage_NoUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	req, _ := http.NewRequest("GET", "/user/login", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User Login")
}

func TestGetLoginPage_WithUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusFound, resp.Code)
	assert.NotContains(t, resp.Body.String(), "User Login")
}

func TestGetProfilePage_NoUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	req, _ := http.NewRequest("GET", "/user/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "User is not logged in.")
}

func TestGetProfilePage_WithUser(t *testing.T) {
	// Setup
	router := setupUserRouter()
	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User name: Test User")
}
