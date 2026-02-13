//go:build L1

package controller_test

import (
	"glower/auth"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"glower/middleware"
	"glower/tests/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Setup

func setupAuthRouter(mockRepo repository.AuthRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r)

	group := r.Group("/auth")
	factory := func(c *gin.Context) repository.AuthRepository { return mockRepo }

	group.POST("/signup", controller.CreateRegister(factory))
	group.POST("/login", controller.CreateLogin(factory))
	group.POST("/logout", middleware.CreateAuth(true), controller.CreateLogout())

	return r
}

// Suite

type authControllerSuite struct {
	suite.Suite
	mockRepo *mocks.AuthRepoMock
	router   *gin.Engine
	token    string
}

func (s *authControllerSuite) SetupSuite() {
	var err error
	s.token, err = mocks.CreateTokenMock()
	s.Require().NoError(err)
}

func (s *authControllerSuite) SetupTest() {
	s.mockRepo = new(mocks.AuthRepoMock)
	s.router = setupAuthRouter(s.mockRepo)
}

func (s *authControllerSuite) newRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(&http.Cookie{Name: auth.AccessTokenName, Value: s.token})
	return req
}

// Tests

func (s *authControllerSuite) TestLogin_NoData() {
	// Arrange
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/auth/login", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func (s *authControllerSuite) TestLogin_UnableToGetUser() {
	// Arrange
	s.mockRepo.On("GetUser", mock.Anything).Return(model.User{}, assert.AnError)

	form := url.Values{}
	form.Add("email", "test@example.com")
	form.Add("password", "Password123")

	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/auth/login", strings.NewReader(form.Encode()))

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusUnauthorized, resp.Code)
	s.Contains(resp.Body.String(), "Invalid email or password.")
}

func (s *authControllerSuite) TestLogin_WithData() {
	// Arrange
	s.mockRepo.On("GetUser", mock.Anything).Return(mocks.GetTestUser(), nil)

	for _, testData := range mocks.GetLoginDataTestCases() {
		form := url.Values{}
		form.Add("email", testData.Email)
		form.Add("password", testData.Password)

		req := s.newRequest("POST", "/auth/login", strings.NewReader(form.Encode()))
		resp := httptest.NewRecorder()

		// Act
		s.router.ServeHTTP(resp, req)

		// Assert
		s.Equal(testData.ExpectedCode, resp.Code)
		s.Contains(resp.Body.String(), testData.ExpectedMsg)
	}
}

func (s *authControllerSuite) TestRegister_NoData() {
	// Arrange
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/auth/signup", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func (s *authControllerSuite) TestRegister_UnableToInsertUser() {
	// Arrange
	s.mockRepo.On("InsertUser", mock.Anything).Return(assert.AnError)

	form := url.Values{}
	form.Add("first-name", "fs")
	form.Add("last-name", "ls")
	form.Add("email", "test@example.com")
	form.Add("password", "Qazxc12345!@#$%")
	form.Add("confirm-password", "Qazxc12345!@#$%")

	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/auth/signup", strings.NewReader(form.Encode()))

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Cannot register user. Please try again later.")
}

func (s *authControllerSuite) TestRegister_WithData() {
	// Arrange
	s.mockRepo.On("InsertUser", mock.Anything).Return(nil)

	for _, testData := range mocks.GetRegisterDataTestCases() {
		form := url.Values{}
		form.Add("first-name", testData.FirstName)
		form.Add("last-name", testData.LastName)
		form.Add("email", testData.Email)
		form.Add("password", testData.Password)
		form.Add("confirm-password", testData.ConfirmPassword)

		req := s.newRequest("POST", "/auth/signup", strings.NewReader(form.Encode()))
		resp := httptest.NewRecorder()

		// Act
		s.router.ServeHTTP(resp, req)

		// Assert
		s.Equal(testData.ExpectedCode, resp.Code)
		s.Contains(resp.Body.String(), testData.ExpectedMsg)
	}
}

func (s *authControllerSuite) TestLogout_WithToken() {
	// Arrange
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/auth/logout", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Equal("/?oper=logout", resp.Header().Get("HX-Redirect"))
	for _, cookie := range resp.Result().Cookies() {
		if cookie.Name == auth.AccessTokenName {
			s.True(cookie.Value == "" && cookie.MaxAge < 0)
		}
		if cookie.Name == auth.RefreshTokenName {
			s.True(cookie.Value == "" && cookie.MaxAge < 0)
		}
	}
}

func (s *authControllerSuite) TestLogout_NoToken() {
	// Arrange
	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/auth/logout", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusUnauthorized, resp.Code)
	s.Contains(resp.Body.String(), "User is not logged in.")
}

func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(authControllerSuite))
}
