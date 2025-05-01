package controller

import (
	"database/sql"
	"glower/auth"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"glower/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock

type mockCartRepo struct{}

func (r *mockCartRepo) GetUserCart(userId uint) (model.Cart, error) {
	return model.Cart{}, nil
}

func (r *mockCartRepo) AddOrUpdateCartItem(cartID, flowerID uint) (model.CartItem, error) {
	return model.CartItem{}, nil
}

func (r *mockCartRepo) GetCartItems(cartID uint) ([]model.CartItem, error) {
	return []model.CartItem{
		{
			CartID:   1,
			FlowerID: 11,
			Flower: model.Flower{
				Name:        "Sunflower",
				Price:       9.99000,
				Available:   false,
				Description: "Yellow flower",
				Inventory: model.Inventory{
					FlowerID: 1,
					Stock:    10,
				},
			},
			Quantity: 1,
		},
		{
			CartID:   2,
			FlowerID: 12,
			Flower: model.Flower{
				Name:          "Poppy",
				Price:         7.99000,
				Available:     true,
				Description:   "Red flower",
				DiscountPrice: sql.NullFloat64{Float64: 5.99000, Valid: true},
				Inventory: model.Inventory{
					FlowerID: 1,
					Stock:    13,
				},
			},
			Quantity: 2,
		},
	}, nil
}

func (r *mockCartRepo) RemoveCartItem(cartID uint, cartItemID uint) error {
	return nil
}

func (r *mockCartRepo) GetFlowerByID(flowerID uint) (model.Flower, error) {
	return model.Flower{}, nil
}

func GetMockCartRepoFactory() repository.CartRepoFactory {
	return func(c *gin.Context) repository.CartRepository {
		return &mockCartRepo{}
	}
}

// Setup

func setupCartRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/cart")
	factory := GetMockCartRepoFactory()
	group.GET("/", middleware.CreateAuth(true), controller.CreateGetCartItems(factory))
	group.POST("/", middleware.CreateAuth(true), controller.CreateAddCartItem(factory))
	group.DELETE("/:id", middleware.CreateAuth(true), controller.CreateRemoveCartItem(factory))

	return r
}

// Tests

func TestGetCartItems_MapsDataCorrectly(t *testing.T) {
	// Arrange
	router := setupCartRouter()

	token, err := createTokenMock()
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", "/cart/", nil)
	req.AddCookie(&http.Cookie{
		Name:  auth.AccessTokenName,
		Value: token,
	})
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	assert.Contains(t, resp.Body.String(), "Sunflower")
	assert.Contains(t, resp.Body.String(), "$9.99")
	assert.Contains(t, resp.Body.String(), "<td>1</td>")

	assert.Contains(t, resp.Body.String(), "Poppy")
	assert.Contains(t, resp.Body.String(), "$5.99")
	assert.Contains(t, resp.Body.String(), "<td>2</td>")

	assert.NotContains(t, resp.Body.String(), "Your cart is empty.")
	assert.Contains(t, resp.Body.String(), "Total Price: $21.97")
}
