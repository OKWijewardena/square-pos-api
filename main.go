package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/OKWijewardena/square-pos-api/auth"
	"github.com/OKWijewardena/square-pos-api/database"
	"github.com/OKWijewardena/square-pos-api/handlers"
	"github.com/OKWijewardena/square-pos-api/models"
	"github.com/OKWijewardena/square-pos-api/square"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

type CreateOrderRequest struct {
	Table string             `json:"table"`
	Items []square.OrderItem `json:"items"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	log.Println("✅ .env file loaded")
	log.Println("🧪 Loaded Square Access Token:", os.Getenv("SQUARE_ACCESS_TOKEN"))
	log.Println("🧪 Loaded Square Location ID:", os.Getenv("SQUARE_LOCATION_ID"))

	token := os.Getenv("SQUARE_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("❌ SQUARE_ACCESS_TOKEN is empty")
	} else {
		log.Println("🔑 Square Token Loaded (partial):", token[:10], "...")
	}

	database.Connect()
	database.DB.AutoMigrate(
		&models.Restaurant{},
		&models.Order{},
		&models.OrderItem{},
		&models.Modifier{},
		&models.Discount{},
		&models.Payment{},
	)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(RequestIDMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		requestID := c.GetString("RequestID")
		c.JSON(200, gin.H{"message": "pong", "request_id": requestID})
	})

	router.POST("/login", handlers.Login)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/dashboard", func(c *gin.Context) {
			restaurantID := c.GetUint("restaurant_id")
			c.JSON(200, gin.H{"message": "Welcome!", "restaurant_id": restaurantID})
		})

		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders/:order_id", handlers.GetOrderByID)
		protected.GET("/orders/table/:table_number", handlers.GetOrdersByTable)
		protected.POST("/payments", handlers.SubmitPayment)
	}

	router.POST("/square/order", func(c *gin.Context) {
		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		response, err := square.CreateSquareOrder(req.Table, req.Items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "detail": err.Error()})
			return
		}

		var squareResp map[string]interface{}
		if err := json.Unmarshal(response, &squareResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read Square response", "detail": err.Error()})
			return
		}

		orderRaw, ok := squareResp["order"]
if !ok {
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "Square response does not include order data",
		"detail": squareResp,
	})
	return
}

orderMap, ok := orderRaw.(map[string]interface{})
if !ok {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Invalid format for Square order",
	})
	return
}

customResp := square.TransformSquareOrder(map[string]interface{}{"order": orderMap})
c.JSON(http.StatusOK, customResp)

	})

	router.Run(":8080")
}