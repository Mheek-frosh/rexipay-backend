package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your_secret_key" // Fallback for dev if .env fails
	}
	return []byte(secret)
}

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTKey())
}
