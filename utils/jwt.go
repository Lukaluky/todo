package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte(getEnv("JWT_SECRET","something"))

type Claims struct {

	UserID  uint `json:"user_id"`
	jwt.RegisteredClaims

}

func GenerateJWT(userID uint) (string, error){

    ttlHours, _ := strconv.Atoi(getEnv("JWT_TTL_HOURS","24"))
	exp := time.Now().Add(time.Duration(ttlHours)*time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(secret)

}

func ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}
	return claims, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

