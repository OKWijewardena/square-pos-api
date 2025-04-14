package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/OKWijewardena/square-pos-api/auth"
	"github.com/OKWijewardena/square-pos-api/database"
	"github.com/OKWijewardena/square-pos-api/models"
)

type LoginRequest struct {
	RestaurantID uint `json:"restaurant_id"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Optional: validate if restaurant exists
	var restaurant models.Restaurant
	if err := database.DB.First(&restaurant, req.RestaurantID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Restaurant not found"})
		return
	}

	token, err := auth.GenerateJWT(req.RestaurantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
