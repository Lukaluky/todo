package middlewares

import (
	
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	
	"todo/utils"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        raw := c.GetHeader("Authorization")
        if raw == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }

        // Bearer xxxx
        token := strings.TrimPrefix(raw, "Bearer ")

        // Парсим JWT
        claims, err := utils.ParseJWT(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            return
        }

        // Дальше НИЧЕГО НЕ ПРОВЕРЯЕМ В REDIS !!!

        c.Set("userID", claims.UserID)
        c.Next()
    }
}
