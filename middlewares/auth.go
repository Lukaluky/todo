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
		token := strings.TrimPrefix(raw,"Bearer ")

		claims, err := utils.ParseJWT(token)
        if err != nil {
	         c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	        return
        }

		c.Set("userID", claims.UserID)
	    c.Next()
	}
}