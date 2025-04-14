package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/OKWijewardena/square-pos-api/database"
	"github.com/OKWijewardena/square-pos-api/models"
)

func CreateOrder(c *gin.Context) {
	var req struct {
		Table string `json:"table"`
		Items []struct {
			Name      string  `json:"name"`
			UnitPrice float64 `json:"unit_price"`
			Quantity  int     `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Table == "" || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: table and items are required"})
		return
	}

	restaurantID := c.GetUint("restaurant_id")

	order := models.Order{
		RestaurantID: restaurantID,
		TableNumber:  req.Table,
		OpenedAt:     time.Now(),
		IsClosed:     false,
	}

	for _, item := range req.Items {
		if item.Name == "" || item.Quantity <= 0 || item.UnitPrice < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item data"})
			return
		}
		order.OrderItems = append(order.OrderItems, models.OrderItem{
			Name:      item.Name,
			UnitPrice: item.UnitPrice,
			Quantity:  item.Quantity,
			Amount:    item.UnitPrice * float64(item.Quantity),
		})
	}

	database.DB.Create(&order)

	c.JSON(http.StatusOK, gin.H{"message": "Order created", "order_id": order.ID})
}

func GetOrderByID(c *gin.Context) {
	restaurantID := c.GetUint("restaurant_id")
	orderID := c.Param("order_id")

	var order models.Order
	err := database.DB.Preload("OrderItems").
		Where("id = ? AND restaurant_id = ?", orderID, restaurantID).
		First(&order).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetOrdersByTable(c *gin.Context) {
	table := c.Param("table_number")
	restaurantID := c.GetUint("restaurant_id")

	var orders []models.Order
	database.DB.Where("restaurant_id = ? AND table_number = ?", restaurantID, table).
		Preload("OrderItems").Find(&orders)

	c.JSON(http.StatusOK, orders)
}

func SubmitPayment(c *gin.Context) {
	var req struct {
		OrderID    uint    `json:"order_id"`
		BillAmount float64 `json:"billAmount"`
		TipAmount  float64 `json:"tipAmount"`
		PaymentID  string  `json:"paymentId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.OrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	restaurantID := c.GetUint("restaurant_id")

	// Check if the order belongs to this restaurant
	var order models.Order
	err := database.DB.Where("id = ? AND restaurant_id = ?", req.OrderID, restaurantID).First(&order).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found for your restaurant"})
		return
	}

	payment := models.Payment{
		OrderID:    req.OrderID,
		BillAmount: req.BillAmount,
		TipAmount:  req.TipAmount,
		PaymentID:  req.PaymentID,
	}
	database.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{"message": "Payment submitted"})
}
