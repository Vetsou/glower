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
	"github.com/stretchr/testify/suite"
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

// Suite

type userControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	token  string
}

func (s *userControllerTestSuite) SetupSuite() {
	s.router = setupUserRouter()

	var err error
	s.token, err = createTokenMock()
	s.Require().NoError(err)
}

// Tests

func (s *userControllerTestSuite) TestGetRegisterPage_NoUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/register", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "User Register")
}

func (s *userControllerTestSuite) TestGetRegisterPage_WithUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/register", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: s.token,
	})
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusFound, resp.Code)
	s.NotContains(resp.Body.String(), "User Register")
}

func (s *userControllerTestSuite) TestGetLoginPage_NoUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/login", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "User Login")
}

func (s *userControllerTestSuite) TestGetLoginPage_WithUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/login", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: s.token,
	})
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusFound, resp.Code)
	s.NotContains(resp.Body.String(), "User Login")
}

func (s *userControllerTestSuite) TestGetProfilePage_NoUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/profile", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusUnauthorized, resp.Code)
	s.Contains(resp.Body.String(), "User is not logged in.")
}

func (s *userControllerTestSuite) TestGetProfilePage_WithUser() {
	// Arrange
	req, _ := http.NewRequest("GET", "/user/profile", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: s.token,
	})
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "User name: Test User")
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(userControllerTestSuite))
}
