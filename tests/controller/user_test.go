package controller

import (
	"glower/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRegisterPage_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	req, _ := http.NewRequest("GET", "/user/register", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User Register")
}

func TestGetRegisterPage_WithUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	token, err := setupTokenAuth()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/register", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusFound, resp.Code)
	assert.NotContains(t, resp.Body.String(), "User Register")
}

func TestGetLoginPage_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	req, _ := http.NewRequest("GET", "/user/login", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User Login")
}

func TestGetLoginPage_WithUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	token, err := setupTokenAuth()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusFound, resp.Code)
	assert.NotContains(t, resp.Body.String(), "User Login")
}

func TestGetProfilePage_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "User is not logged in.")
}

func TestGetProfilePage_WithUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupUserRouter()

	token, err := setupTokenAuth()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User name: Test User")
}
