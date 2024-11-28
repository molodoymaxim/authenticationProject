package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	UserID      string
	HashedToken string
	AccessToken string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"<access_token>"`
	RefreshToken string `json:"refresh_token" example:"<refresh_token>"`
}

type TokenRequest struct {
	UserID string `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"<refresh_token>"`
}
