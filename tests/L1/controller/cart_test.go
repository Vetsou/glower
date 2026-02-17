//go:build L1

package l1_test

import (
	"glower/auth"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"glower/middleware"
	"glower/tests/L1/mocks"
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

func setupCartRouter(mockRepo repository.CartRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r)

	group := r.Group("/cart", middleware.CreateAuth(true))
	factory := func(c *gin.Context) repository.CartRepository { return mockRepo }

	group.GET("/", controller.CreateGetCartItems(factory))
	group.POST("/", controller.CreateAddCartItem(factory))
	group.DELETE("/:id", controller.CreateRemoveCartItem(factory))

	return r
}

// Suite

type cartControllerSuite struct {
	suite.Suite
	mockRepo *mocks.CartRepoMock
	router   *gin.Engine
	token    string
}

func (s *cartControllerSuite) SetupSuite() {
	var err error
	s.token, err = mocks.CreateTokenMock()
	s.Require().NoError(err)
}

func (s *cartControllerSuite) SetupTest() {
	s.mockRepo = new(mocks.CartRepoMock)
	s.router = setupCartRouter(s.mockRepo)
}

func (s *cartControllerSuite) newRequest(method, path string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(&http.Cookie{Name: auth.AccessTokenName, Value: s.token})
	return req
}

// Tests

func (s *cartControllerSuite) TestGetCartItems_MapsDataCorrectly() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("GetCartItems", mock.Anything).Return(mocks.GetTestCartItems(), nil)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/cart/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	body := resp.Body.String()

	s.Contains(body, "Sunflower")
	s.Contains(body, "$9.99")
	s.Contains(body, "<td>1</td>")

	s.Contains(body, "Poppy")
	s.Contains(body, "$5.99")
	s.Contains(body, "<td>2</td>")

	s.NotContains(body, "Your cart is empty.")
	s.Contains(body, "Total Price: $21.97")
}

func (s *cartControllerSuite) TestGetCartItems_EmptyCart() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("GetCartItems", mock.Anything).Return(mocks.GetEmptyCartItems(), nil)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/cart/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "Your cart is empty.")
}

func (s *cartControllerSuite) TestGetCartItems_UserCartError() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, assert.AnError)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/cart/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load your cart.")
}

func (s *cartControllerSuite) TestGetCartItems_CartItemsError() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("GetCartItems", mock.Anything).Return([]model.CartItem{}, assert.AnError)

	resp := httptest.NewRecorder()
	req := s.newRequest("GET", "/cart/", nil)

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load your cart items.")
}

func (s *cartControllerSuite) TestAddCartItem_WithCorrectData() {
	// Arrange
	s.mockRepo.On("DecreaseInventoryAndGetFlower", mock.Anything).Return(mocks.GetCartFlower(true), nil)
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("AddOrUpdateCartItem", mock.Anything, mock.Anything).Return(mocks.GetAddedCartFlower(), nil)

	form := url.Values{}
	form.Add("flowerId", "13")

	resp := httptest.NewRecorder()
	req := s.newRequest("POST", "/cart/", strings.NewReader(form.Encode()))

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "Flower Sunflower was added to your cart. You currently have 1 Sunflower in your cart.")
}

func (s *cartControllerSuite) TestAddCartItem_InvalidFlowerID() {
	// Arrange
	s.mockRepo.On("DecreaseInventoryAndGetFlower", mock.Anything).Return(model.Flower{}, assert.AnError)

	form := url.Values{}
	form.Add("flowerId", "13")

	req := s.newRequest("POST", "/cart/", strings.NewReader(form.Encode()))
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load flower data.")
}

func (s *cartControllerSuite) TestAddCartItem_FlowerOutOfStock() {
	// Arrange
	s.mockRepo.On("DecreaseInventoryAndGetFlower", mock.Anything).Return(model.Flower{}, repository.ErrOutOfStock)

	form := url.Values{}
	form.Add("flowerId", "13")

	req := s.newRequest("POST", "/cart/", strings.NewReader(form.Encode()))
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "The requested flower is out of stock.")
}

func (s *cartControllerSuite) TestAddCartItem_UnableToLoadCart() {
	// Arrange
	s.mockRepo.On("DecreaseInventoryAndGetFlower", mock.Anything).Return(mocks.GetCartFlower(true), nil)
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, assert.AnError)

	form := url.Values{}
	form.Add("flowerId", "13")

	req := s.newRequest("POST", "/cart/", strings.NewReader(form.Encode()))
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load your cart.")
}

func (s *cartControllerSuite) TestAddCartItem_UnableToAddCartItem() {
	// Arrange
	s.mockRepo.On("DecreaseInventoryAndGetFlower", mock.Anything).Return(mocks.GetCartFlower(true), nil)
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("AddOrUpdateCartItem", mock.Anything, mock.Anything).Return(model.CartItem{}, assert.AnError)

	form := url.Values{}
	form.Add("flowerId", "13")

	req := s.newRequest("POST", "/cart/", strings.NewReader(form.Encode()))
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load your cart items.")
}

func (s *cartControllerSuite) TestRemoveCartItem_WithCorrectData() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("RemoveCartItem", mock.Anything, mock.Anything).Return(nil)

	req := s.newRequest("DELETE", "/cart/1", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusOK, resp.Code)
	s.Contains(resp.Body.String(), "Item was removed from your cart.")
}

func (s *cartControllerSuite) TestRemoveCartItem_WrongId() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("RemoveCartItem", mock.Anything, mock.Anything).Return(nil)

	req := s.newRequest("DELETE", "/cart/test", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusBadRequest, resp.Code)
	s.Contains(resp.Body.String(), "Wrong cart item ID.")
}

func (s *cartControllerSuite) TestRemoveCartItem_UnableToLoadCart() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, assert.AnError)
	s.mockRepo.On("RemoveCartItem", mock.Anything, mock.Anything).Return(nil)

	req := s.newRequest("DELETE", "/cart/1", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to load your cart.")
}

func (s *cartControllerSuite) TestRemoveCartItem_UnableToRemoveCartItem() {
	// Arrange
	s.mockRepo.On("GetUserCart", mock.Anything).Return(model.Cart{}, nil)
	s.mockRepo.On("RemoveCartItem", mock.Anything, mock.Anything).Return(assert.AnError)

	req := s.newRequest("DELETE", "/cart/1", nil)
	resp := httptest.NewRecorder()

	// Act
	s.router.ServeHTTP(resp, req)

	// Assert
	s.Equal(http.StatusInternalServerError, resp.Code)
	s.Contains(resp.Body.String(), "Unable to remove cart item.")
}

func TestCartControllerTestSuite(t *testing.T) {
	suite.Run(t, new(cartControllerSuite))
}
