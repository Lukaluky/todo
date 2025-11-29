package handlers

import (
	"net/http"
	"time"

	"todo/db"
	"todo/models"
	"todo/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const RefreshCookieName = "refresh_token"

// -------------------------- REGISTER ----------------------------

func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{Email: input.Email, Password: string(hash)}

	if err := db.PG.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	// ---- access token ----
	access, _ := utils.GenerateJWT(user.ID)

	// ---- refresh token ----
	refresh, exp, _ := utils.GenerateRefreshToken()

	// ---- Save refresh token in Redis ----
	db.Redis.Set(c, "rt:"+refresh, user.ID, time.Until(exp))

	// ---- Set refresh as HttpOnly Cookie ----
	c.SetCookie(
		RefreshCookieName,
		refresh,
		int(time.Until(exp).Seconds()),
		"/",
		"",
		false,  
		true,   
	)

	c.JSON(200, gin.H{
		"message":      "registered",
		"access_token": access,
	})
}

// -------------------------- LOGIN ----------------------------

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := db.PG.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// ---- Generate tokens ----
	access, _ := utils.GenerateJWT(user.ID)
	refresh, exp, _ := utils.GenerateRefreshToken()

	// ---- Save refresh to Redis ----
	db.Redis.Set(c, "rt:"+refresh, user.ID, time.Until(exp))

	// ---- Set cookie ----
	c.SetCookie(
		RefreshCookieName,
		refresh,
		int(time.Until(exp).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"access_token": access,
	})
}
