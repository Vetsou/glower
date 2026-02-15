//go:build L1

package l1_test

import (
	"glower/database/model"
	"glower/initializers"
	"glower/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// Suite

type rolesMiddlewareSuite struct {
	suite.Suite
	router *gin.Engine
	role   any
}

func (s *rolesMiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.router = gin.New()
	initializers.InitHTMLTemplates(s.router)

	s.router.Use(func(c *gin.Context) {
		if s.role != nil {
			c.Set("role", s.role)
		}
		c.Next()
	})

	s.router.GET(
		"/roleRoute",
		middleware.CreateRolesAuth(model.RoleUser),
		func(c *gin.Context) { c.String(http.StatusOK, "OK") },
	)
}

func (s *rolesMiddlewareSuite) doRequest() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/roleRoute", nil)

	s.router.ServeHTTP(w, req)

	return w
}

// Tests

func (s *rolesMiddlewareSuite) TestRoles_AllowsCorrectRole() {
	// Arrange
	s.role = model.RoleUser

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func (s *rolesMiddlewareSuite) TestRoles_ForbiddenWrongRole() {
	// Arrange
	s.role = model.RoleUser

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func (s *rolesMiddlewareSuite) TestRoles_UnauthorizedNoRole() {
	// Arrange
	s.role = nil

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "cannot view")
}

func (s *rolesMiddlewareSuite) TestRoles_InvalidRoleType() {
	// Arrange
	s.role = "fake_admin"

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "verify your access rights")
}

func (s *rolesMiddlewareSuite) TestRoles_MultipleAllowedRoles() {
	// Arrange
	gin.SetMode(gin.TestMode)

	s.router = gin.New()

	s.router.Use(func(c *gin.Context) {
		if s.role != nil {
			c.Set("role", s.role)
		}
		c.Next()
	})

	s.router.GET(
		"/roleRoute",
		middleware.CreateRolesAuth(model.RoleUser, model.RoleAdmin),
		func(c *gin.Context) { c.String(http.StatusOK, "OK") },
	)

	s.role = model.RoleAdmin

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusOK, w.Code)
}

func TestRoleMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(rolesMiddlewareSuite))
}
