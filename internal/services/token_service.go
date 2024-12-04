package services

import (
	"authenticationProject/internal/models"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TokenService struct {
	SecretKey   string
	AccessToken time.Duration
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		SecretKey:   secretKey,
		AccessToken: time.Hour * 1,
	}
}

func (s *TokenService) GenerateAccessToken(userID, pairID, clientIP string) (string, error) {
	claims := &models.CustomClaims{
		UserID: userID,
		IP:     clientIP,
		PairID: pairID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.AccessToken)),
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

func (s *TokenService) GenerateRefreshToken(userID, pairID, clientIP string) (string, string, error) {
	refreshTokenPlain := fmt.Sprintf("%s:%s:%s", userID, pairID, clientIP)
	refreshTokenEncoded := base64.StdEncoding.EncodeToString([]byte(refreshTokenPlain))
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenPlain), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return refreshTokenEncoded, string(hashedToken), nil
}
