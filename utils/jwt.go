package utils

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret")

func GenerateTokens(userID, ip string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken := "some-unique-refresh-token" // Use a unique generation mechanism

	return accessTokenString, refreshToken, nil
}

func GetValidator() *validator.Validate {
	return validator.New()
}
