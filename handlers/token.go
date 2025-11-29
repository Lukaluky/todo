package handlers

import (
	
	"net/http"
	"time"
	"todo/db"
	"todo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Refresh — обновление access token по refresh cookie.
func Refresh(c *gin.Context) {
	refresh, err := c.Cookie(RefreshCookieName)
	if err != nil || refresh == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh cookie missing"})
		return
	}

	// Ищем refresh token в Redis
	userIDstr, err := db.Redis.Get(db.Ctx, "rt:"+refresh).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}
	uid64, err := strconv.ParseUint(userIDstr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id in redis"})
		return
	}
	userID := uint(uid64)

	// Генерируем новый access
	access, err := utils.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot generate access token"})
		return
	}

	// Ротация refresh token
	newRefresh, exp, _ := utils.GenerateRefreshToken()

	pipe := db.Redis.TxPipeline()
	pipe.Del(db.Ctx, "rt:"+refresh)
	pipe.Set(db.Ctx, "rt:"+newRefresh, userID, time.Until(exp))
	_, err = pipe.Exec(db.Ctx)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
        return
	}
	// Устанавливаем новую refresh cookie
	c.SetCookie(
		RefreshCookieName,
		newRefresh,
		int(time.Until(exp).Seconds()),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
	})
}

// Logout — удаляет refresh cookie и запись в Redis
func Logout(c *gin.Context) {
	refresh, err := c.Cookie(RefreshCookieName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh cookie missing"})
		return
	}

	db.Redis.Del(db.Ctx, "rt:"+refresh)

	// Удаляем cookie (ставим maxAge = -1)
	c.SetCookie(
		RefreshCookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
