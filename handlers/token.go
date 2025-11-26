package handlers

import (
	"net/http"
	"time"
	"todo/db"
	"todo/models"
	"todo/utils"

	"github.com/gin-gonic/gin"
)

// Refresh обновляет access token по refresh token
func Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	var rt models.RefreshToken
	if err := db.PG.Where("token = ? AND expires_at > ?", body.RefreshToken, time.Now()).First(&rt).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	
	access, err := utils.GenerateJWT(rt.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
	})
}

// Logout удаляет refresh token из базы
func Logout(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	result := db.PG.Where("token = ?", body.RefreshToken).Delete(&models.RefreshToken{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete token"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
