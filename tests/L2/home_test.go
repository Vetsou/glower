//go:build L2

package l2_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage_NoOper(t *testing.T) {
	// Act
	resp, err := http.Get(baseURL + "/")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotContains(t, string(body), "logged")
}

func TestHomePage_Logout(t *testing.T) {
	// Act
	resp, err := http.Get(baseURL + "/?oper=logout")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "You have successfully logged out.")
}

func TestHomePage_Login(t *testing.T) {
	// Act
	resp, err := http.Get(baseURL + "/?oper=login")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), "You have successfully logged in.")
}

func TestHomePage_Register(t *testing.T) {
	// Act
	resp, err := http.Get(baseURL + "/?oper=register")
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body),
		"Registration successful. Please log in to continue.",
	)
}
