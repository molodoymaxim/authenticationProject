package services

import (
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateAccessToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	userID := "test-user-id"
	pairID := uuid.NewString()
	clientIP := "127.0.0.1"

	token, err := tokenService.GenerateAccessToken(userID, pairID, clientIP)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	userID := "test-user-id"
	pairID := uuid.NewString()
	clientIP := "127.0.0.1"

	token, err := tokenService.GenerateAccessToken(userID, pairID, clientIP)
	assert.NoError(t, err)

	claims, err := tokenService.ValidateAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, pairID, claims.PairID)
	assert.Equal(t, clientIP, claims.IP)
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	userID := "test-user-id"
	pairID := uuid.NewString()
	clientIP := "127.0.0.1"

	refreshTokenEncoded, hashedToken, err := tokenService.GenerateRefreshToken(userID, pairID, clientIP)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshTokenEncoded)
	assert.NotEmpty(t, hashedToken)

	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshTokenEncoded)
	assert.NoError(t, err)
	refreshTokenPlain := string(refreshTokenBytes)

	err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(refreshTokenPlain))
	assert.NoError(t, err)
}
