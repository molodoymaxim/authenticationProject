package services

import (
	"authenticationProject/internal/models"
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	SecretKey string
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		SecretKey: secretKey,
	}
}

func (s *TokenService) GenerateAccessToken(userID, clientIP string) (string, error) {
	claims := &models.CustomClaims{
		UserID: userID,
		IP:     clientIP,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.SecretKey))
}

func (s *TokenService) ValidateAccessToken(tokenStr string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if claims, ok := token.Claims.(*models.CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s *TokenService) GenerateRefreshToken() (string, string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", "", err
	}
	refreshToken := base64.URLEncoding.EncodeToString(tokenBytes)
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return refreshToken, string(hashedToken), nil
}
