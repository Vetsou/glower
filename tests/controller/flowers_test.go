package controller

import (
	"database/sql"
	"fmt"
	"glower/controller"
	"glower/database/model"
	"glower/database/repository"
	"glower/initializers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock

type mockStockRepo struct{}

func (r *mockStockRepo) AddFlower(flower model.Flower, count uint) error { return nil }

func (r *mockStockRepo) RemoveFlower(id string) error {
	// Simulate DB failing to find flower by ID
	_, err := strconv.Atoi(id)
	return err
}

func (r *mockStockRepo) GetFlowers() ([]model.Flower, error) {
	var flowers = []model.Flower{
		{
			Name:        "Sunflower",
			Price:       9.99,
			Available:   false,
			Description: "Yellow flower",
			Inventory: model.Inventory{
				FlowerID: 1,
				Stock:    10,
			},
		},
		{
			Name:          "Poppy",
			Price:         7.99,
			Available:     true,
			Description:   "Red flower",
			DiscountPrice: sql.NullFloat64{Float64: 5.99, Valid: true},
			Inventory: model.Inventory{
				FlowerID: 1,
				Stock:    13,
			},
		},
	}

	return flowers, nil
}

func GetMockStockRepoFactory() repository.StockRepoFactory {
	return func(c *gin.Context) repository.StockRepository {
		return &mockStockRepo{}
	}
}

// Setup

func setupFlowersRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	initializers.InitHTMLTemplates(r, "../../")

	group := r.Group("/flowers")
	factory := GetMockStockRepoFactory()
	group.GET("/", controller.CreateGetFlowers(factory))
	group.POST("/", controller.CreateAddFlower(factory))
	group.DELETE("/:id", controller.CreateRemoveFlower(factory))

	return r
}

// Test Cases

type deleteTestData struct {
	deleteFlowerId string
	expectedCode   int
	expectedMsg    string
}

var deleteTestCases = []deleteTestData{
	{"1", 200, ""},
	{"", 404, "404 page not found"},
	{"wrongId", 500, "Error deleting flower. Please try again later."},
}

// Tests

func TestGetFlowers_MapsDataCorrectly(t *testing.T) {
	// Arrange
	router := setupFlowersRouter()

	req, _ := http.NewRequest("GET", "/flowers/", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	assert.Contains(t, resp.Body.String(), "Sunflower")
	assert.Contains(t, resp.Body.String(), "$9.99")
	assert.Contains(t, resp.Body.String(), "Description: Yellow flower")
	assert.Contains(t, resp.Body.String(), "Stock: 10")
	assert.Contains(t, resp.Body.String(), "Available: No")

	assert.Contains(t, resp.Body.String(), "Poppy")
	assert.Contains(t, resp.Body.String(), `"text-decoration: line-through;">$7.99<`)
	assert.Contains(t, resp.Body.String(), "$5.99")
	assert.Contains(t, resp.Body.String(), "Description: Red flower")
	assert.Contains(t, resp.Body.String(), "Stock: 13")
	assert.Contains(t, resp.Body.String(), "Available: Yes")
}

func TestRemoveFlower(t *testing.T) {
	// Arrange
	router := setupFlowersRouter()

	for _, testData := range deleteTestCases {
		req, _ := http.NewRequest("DELETE", "/flowers/"+testData.deleteFlowerId, nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, testData.expectedCode, resp.Code)
		assert.Contains(t, resp.Body.String(), testData.expectedMsg)
	}
}

func TestCreateFlower_NoData(t *testing.T) {
	// Arrange
	router := setupFlowersRouter()

	req, _ := http.NewRequest("POST", "/flowers/", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid form data. Please fill all required fields.")
}

func TestCreateFlower_WithCorrectData(t *testing.T) {
	// Arrange
	router := setupFlowersRouter()

	form := url.Values{}
	form.Add("name", "FlowerName")
	form.Add("price", fmt.Sprintf("%f", 15.00))
	form.Add("available", strconv.FormatBool(true))
	form.Add("description", "Nice flower")
	form.Add("discount", fmt.Sprintf("%f", 10.0))
	form.Add("stock", strconv.FormatUint(uint64(12), 10))

	req, _ := http.NewRequest("POST", "/flowers/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "FlowerName")
	assert.Contains(t, resp.Body.String(), `"text-decoration: line-through;">$15<`)
	assert.Contains(t, resp.Body.String(), "$10")
	assert.Contains(t, resp.Body.String(), "Description: Nice flower")
	assert.Contains(t, resp.Body.String(), "Stock: 12")
	assert.Contains(t, resp.Body.String(), "Available: Yes")
}

func TestCreateFlower_WithMissingOptionalData(t *testing.T) {
	// Arrange
	router := setupFlowersRouter()

	form := url.Values{}
	form.Add("name", "FlowerName")
	form.Add("price", fmt.Sprintf("%f", 12.00))
	form.Add("stock", strconv.FormatUint(uint64(13), 10))

	req, _ := http.NewRequest("POST", "/flowers/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "FlowerName")
	assert.Contains(t, resp.Body.String(), "$12")
	assert.Contains(t, resp.Body.String(), "Stock: 13")
	assert.Contains(t, resp.Body.String(), "Available: No")

	assert.NotContains(t, resp.Body.String(), "Description:")
}
