package middleware

import (
	"glower/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Structs

type TestModel struct {
	ID   uint
	Name string
}

// Suite

type transactionMiddlewareSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
}

func (s *transactionMiddlewareSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	// In-memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	err = db.AutoMigrate(&TestModel{})
	s.Require().NoError(err)

	s.db = db
	s.router = gin.New()
	s.router.Use(gin.Recovery())
}

func (s *transactionMiddlewareSuite) createRoute(
	status int,
	handler gin.HandlerFunc,
) {
	s.router.GET(
		"/tx",
		middleware.CreateTransaction(s.db),
		func(c *gin.Context) {
			handler(c)
			c.Status(status)
		},
	)
}

func (s *transactionMiddlewareSuite) doRequest() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tx", nil)

	s.router.ServeHTTP(w, req)

	return w
}

func (s *transactionMiddlewareSuite) countRows() int64 {
	var count int64
	s.db.Model(&TestModel{}).Count(&count)
	return count
}

// Tests

func (s *transactionMiddlewareSuite) TestTransaction_CommitsOnOK() {
	// Arrange
	s.createRoute(http.StatusOK, func(c *gin.Context) {
		tx := c.MustGet("tx").(*gorm.DB)

		err := tx.Create(&TestModel{
			Name: "test",
		}).Error

		s.Require().NoError(err)
	})

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.Equal(int64(1), s.countRows())
}

func (s *transactionMiddlewareSuite) TestTransaction_CommitsOnCreated() {
	// Arrange
	s.createRoute(http.StatusCreated, func(c *gin.Context) {
		tx := c.MustGet("tx").(*gorm.DB)

		err := tx.Create(&TestModel{
			Name: "created",
		}).Error

		s.Require().NoError(err)
	})

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusCreated, w.Code)
	s.Equal(int64(1), s.countRows())
}

func (s *transactionMiddlewareSuite) TestTransaction_RollbackOnBadRequest() {
	// Arrange
	s.createRoute(http.StatusBadRequest, func(c *gin.Context) {
		tx := c.MustGet("tx").(*gorm.DB)

		err := tx.Create(&TestModel{
			Name: "should rollback",
		}).Error

		s.Require().NoError(err)
	})

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal(int64(0), s.countRows())
}

func (s *transactionMiddlewareSuite) TestTransaction_RollbackOnPanic() {
	// Arrange
	s.createRoute(http.StatusOK, func(c *gin.Context) {
		tx := c.MustGet("tx").(*gorm.DB)

		err := tx.Create(&TestModel{
			Name: "panic",
		}).Error

		s.Require().NoError(err)

		panic("something went wrong")
	})

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusInternalServerError, w.Code)
	s.Equal(int64(0), s.countRows())
}

func (s *transactionMiddlewareSuite) TestTransaction_SetsTxInContext() {
	// Arrange
	var txExists bool

	s.createRoute(http.StatusOK, func(c *gin.Context) {
		_, ok := c.Get("tx")
		txExists = ok
	})

	// Act
	w := s.doRequest()

	// Assert
	s.Equal(http.StatusOK, w.Code)
	s.True(txExists)
}

func TestTransactionMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(transactionMiddlewareSuite))
}
