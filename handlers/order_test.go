package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/OKWijewardena/square-pos-api/database"
	"github.com/OKWijewardena/square-pos-api/models"
)

func init() {
	database.Connect()
}

// setupRouter sets up a test Gin engine with required routes.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// Fake JWT middleware for testing
	r.Use(func(c *gin.Context) {
		c.Set("restaurant_id", uint(1)) // simulate authenticated restaurant
		c.Next()
	})

	r.POST("/api/orders", CreateOrder)
	r.POST("/api/payments", SubmitPayment)

	return r
}

func TestCreateOrderHandler(t *testing.T) {
	router := setupRouter()

	// Insert fake restaurant
	database.DB.Create(&models.Restaurant{
		Name: "Test Restaurant",
	})

	// Prepare test request
	payload := map[string]interface{}{
		"table": "5",
		"items": []map[string]interface{}{
			{
				"name":       "Pizza",
				"unit_price": 1000,
				"quantity":   2,
			},
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	if !bytes.Contains(resp.Body.Bytes(), []byte("order_id")) {
		t.Errorf("Expected response to contain 'order_id'")
	}
}

func TestSubmitPaymentHandler(t *testing.T) {
	router := setupRouter()

	// Create dummy order
	order := models.Order{
		RestaurantID: 1,
		TableNumber:  "5",
	}
	database.DB.Create(&order)

	payload := map[string]interface{}{
		"order_id":    order.ID,
		"billAmount":  2000,
		"tipAmount":   300,
		"paymentId":   "test-payment-123",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/api/payments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	if !bytes.Contains(resp.Body.Bytes(), []byte("Payment submitted")) {
		t.Errorf("Expected response to confirm payment submission")
	}
}

