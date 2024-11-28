package models

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims представляет клеймы JWT токена
type CustomClaims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.RegisteredClaims
}

// RefreshTokenData содержит данные о Refresh токене в базе данных
type RefreshTokenData struct {
	UserID      string
	HashedToken string
	AccessToken string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
