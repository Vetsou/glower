package controller

import (
	"glower/controller"
	"glower/initializers"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
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

// Suite

type homeControllerSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *homeControllerSuite) SetupSuite() {
	s.router = setupHomeRouter()
}

// Tests

func (s *homeControllerSuite) TestGetHomePage_NoOper() {
	// Arrange
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	body := resp.Body.String()

	s.NotContains(body, "logged out.")
	s.NotContains(body, "logged in.")
	s.NotContains(body, "Registration successful.")
}

func (s *homeControllerSuite) TestGetHomePage_Logout() {
	// Arrange
	req, _ := http.NewRequest("GET", "/?oper=logout", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "You have successfully logged out.")
}

func (s *homeControllerSuite) TestGetHomePage_Login() {
	// Arrange
	req, _ := http.NewRequest("GET", "/?oper=login", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "You have successfully logged in.")
}

func (s *homeControllerSuite) TestGetHomePage_Register() {
	// Arrange
	req, _ := http.NewRequest("GET", "/?oper=register", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "Registration successful. Please log in to continue.")
}

func TestHomeControllerTestSuite(t *testing.T) {
	suite.Run(t, new(homeControllerSuite))
}
