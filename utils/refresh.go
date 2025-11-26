package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func GenerateRefreshToken() (string, time.Time, error) {
    b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		return "", time.Time{}, err
	}

	refreshtoken := base64.URLEncoding.EncodeToString(b)

	exp := time.Now().Add(7 * 24 * time.Hour)

	return refreshtoken, exp, nil




}