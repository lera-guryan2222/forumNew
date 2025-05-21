package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenManager interface {
	GenerateAccessToken(userID uint) (string, error)
	GenerateRefreshToken() (string, time.Time, error) // Изменено
}

type tokenManager struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration // Добавлено новое поле
}

func NewTokenManager(
	accessSecret,
	refreshSecret string,
	accessExpiry,
	refreshExpiry time.Duration, // Добавлен параметр
) TokenManager {
	return &tokenManager{
		accessTokenSecret:  accessSecret,
		refreshTokenSecret: refreshSecret,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (tm *tokenManager) GenerateRefreshToken() (string, time.Time, error) {
	expiresAt := time.Now().Add(tm.refreshTokenExpiry)
	claims := jwt.MapClaims{
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tm.refreshTokenSecret))
	return tokenString, expiresAt, err
}

// GenerateAccessToken остаётся без изменений

func (tm *tokenManager) GenerateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID, // Теперь принимает uint напрямую
		"exp":     time.Now().Add(tm.accessTokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.accessTokenSecret))
}
