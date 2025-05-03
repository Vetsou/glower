package controller

import (
	"glower/auth"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"glower/tests/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Setup

func setupFlowersRouter(mockRepo repository.StockRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/flowers")
	factory := func(c *gin.Context) repository.StockRepository { return mockRepo }
	group.GET("/", controller.CreateGetFlowers(factory))
	group.POST("/", controller.CreateAddFlower(factory))
	group.DELETE("/:id", controller.CreateRemoveFlower(factory))

	return r
}

// Suite

type stockControllerTestSuite struct {
	suite.Suite
	mockRepo *mocks.StockRepoMock
	router   *gin.Engine
	token    string
}

func (s *stockControllerTestSuite) SetupSuite() {
	var err error
	s.token, err = createTokenMock()
	s.Require().NoError(err)
}

func (s *stockControllerTestSuite) SetupTest() {
	s.mockRepo = new(mocks.StockRepoMock)
	s.router = setupFlowersRouter(s.mockRepo)
}

func (s *stockControllerTestSuite) newRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(&http.Cookie{Name: auth.AccessTokenName, Value: s.token})
	return req
}

// Tests

func (s *stockControllerTestSuite) TestGetFlowers_MapsDataCorrectly() {
	// Arrange
	s.mockRepo.On("GetFlowers", mock.Anything).Return(mocks.GetFlowers(), nil)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/flowers/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	body := resp.Body.String()

	s.Contains(body, "Sunflower")
	s.Contains(body, "$9.99")
	s.Contains(body, "Description: Yellow flower")
	s.Contains(body, "Stock: 10")
	s.Contains(body, "Available: No")

	s.Contains(body, "Poppy")
	s.Contains(body, `"text-decoration: line-through;">$7.99<`)
	s.Contains(body, "$5.99")
	s.Contains(body, "Description: Red flower")
	s.Contains(body, "Stock: 13")
	s.Contains(body, "Available: Yes")
}

func (s *stockControllerTestSuite) TestGetFlowers_UnableToGetFlowers() {
	// Arrange
	s.mockRepo.On("GetFlowers", mock.Anything).Return([]model.Flower{}, assert.AnError)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/flowers/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Failed to load products. Please try again later.")
}

func (s *stockControllerTestSuite) TestRemoveFlower_WrongId() {
	// Arrange
	resp := httptest.NewRecorder()
	req := s.newRequest("DELETE", "/flowers/asd", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "Wrong flower ID.")
}

func (s *stockControllerTestSuite) TestRemoveFlower_UnableToRemoveFlower() {
	// Arrange
	s.mockRepo.On("RemoveFlower", mock.Anything).Return(assert.AnError)

	resp := httptest.NewRecorder()
	req := s.newRequest("DELETE", "/flowers/1", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Error deleting flower. Please try again later.")
}

func (s *stockControllerTestSuite) TestRemoveFlower_ValidData() {
	// Arrange
	s.mockRepo.On("RemoveFlower", mock.Anything).Return(nil)

	resp := httptest.NewRecorder()
	req := s.newRequest("DELETE", "/flowers/1", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
}

func (s *stockControllerTestSuite) TestAddFlower_MapsDataCorrectly() {
	// Arrange
	s.mockRepo.On("AddFlower", mock.Anything, mock.Anything).Return(nil)

	form := mocks.GetValidAddFlowerForm()
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/flowers/", strings.NewReader(form.Encode()))

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	body := resp.Body.String()

	s.Contains(body, "FlowerName")
	s.Contains(body, `"text-decoration: line-through;">$15<`)
	s.Contains(body, "$10")
	s.Contains(body, "Description: Nice flower")
	s.Contains(body, "Stock: 12")
	s.Contains(body, "Available: Yes")
}

func (s *stockControllerTestSuite) TestAddFlower_WithMissingOptionalData() {
	// Arrange
	s.mockRepo.On("AddFlower", mock.Anything, mock.Anything).Return(nil)

	form := mocks.GetOptionalFieldsAddFlowerForm()
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/flowers/", strings.NewReader(form.Encode()))

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	body := resp.Body.String()

	s.Contains(body, "FlowerName")
	s.Contains(body, "$12")
	s.Contains(body, "Stock: 13")
	s.Contains(body, "Available: No")
}

func (s *stockControllerTestSuite) TestAddFlower_NoData() {
	// Arrange
	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/flowers/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func TestStockControllerTestSuite(t *testing.T) {
	suite.Run(t, new(stockControllerTestSuite))
}
