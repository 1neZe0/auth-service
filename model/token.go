// model/token.go
package model

import (
	"errors"
	"time"

	"auth-service/db"
	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	UserID    string `gorm:"primaryKey"`
	TokenHash string
	IP        string
	CreatedAt time.Time
}

func SaveRefreshToken(userID, token, ip string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		TokenHash: string(hash),
		IP:        ip,
		CreatedAt: time.Now(),
	}
	return db.DB.Save(&refreshToken).Error
}

func ValidateRefreshToken(token string) (string, string, error) {
	var refreshToken RefreshToken

	if err := db.DB.First(&refreshToken).Error; err != nil {
		return "", "", errors.New("token not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(refreshToken.TokenHash), []byte(token)); err != nil {
		return "", "", errors.New("invalid token")
	}

	return refreshToken.UserID, refreshToken.IP, nil
}
