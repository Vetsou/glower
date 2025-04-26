package controller

import (
	"glower/auth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLogin_NoLoginData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupAuthRouter()

	token, err := setupTokenAuth()
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/auth/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func TestLogout_WithToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupAuthRouter()

	token, err := setupTokenAuth()
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

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
	gin.SetMode(gin.TestMode)
	router := setupAuthRouter()

	req, _ := http.NewRequest("POST", "/auth/logout", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "User is not logged in.")
}
