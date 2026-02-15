//go:build L1

package l1_test

import (
	"glower/auth"
	"glower/initializers"
	"glower/middleware"
	"glower/tests/L1/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// Setup

func createTestRoute(s *authMiddlewareSuite, isStrict bool) {
	s.router.GET(
		"/authRoute",
		middleware.CreateAuth(isStrict),
		func(c *gin.Context) { c.String(http.StatusOK, "OK") },
	)
}

// Suite

type authMiddlewareSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *authMiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.router = gin.New()
	initializers.InitHTMLTemplates(s.router)
}

func (s *authMiddlewareSuite) doRequest(token *string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/authRoute", nil)

	if token != nil {
		req.AddCookie(&http.Cookie{Name: auth.AccessTokenName, Value: *token})
	}

	s.router.ServeHTTP(w, req)

	return w
}

// Tests

func (s *authMiddlewareSuite) TestAuth_AllowsValidToken() {
	// Arrange
	createTestRoute(s, true)
	token, err := mocks.CreateTokenMock()
	s.Require().NoError(err)

	// Act
	w := s.doRequest(&token)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func (s *authMiddlewareSuite) TestAuth_Strict_RejectsMissingToken() {
	// Arrange
	createTestRoute(s, true)

	// Act
	w := s.doRequest(nil)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "User is not logged in.")
}

func (s *authMiddlewareSuite) TestAuth_NonStrict_AllowsMissingToken() {
	// Arrange
	createTestRoute(s, false)

	// Act
	w := s.doRequest(nil)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func (s *authMiddlewareSuite) TestAuth_Strict_RejectsInvalidToken() {
	// Arrange
	createTestRoute(s, true)
	invalidToken := "invalidtoken123"

	// Act
	w := s.doRequest(&invalidToken)

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "Invalid user credentials.")
}

func (s *authMiddlewareSuite) TestAuth_NonStrict_AllowesInvalidToken() {
	// Arrange
	createTestRoute(s, false)
	invalidToken := "invalidtoken123"

	// Act
	w := s.doRequest(&invalidToken)

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(authMiddlewareSuite))
}
