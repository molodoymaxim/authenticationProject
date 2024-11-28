package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	token, err := tokenService.GenerateAccessToken("test-user-id", "127.0.0.1")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	token, _ := tokenService.GenerateAccessToken("test-user-id", "127.0.0.1")

	claims, err := tokenService.ValidateAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "test-user-id", claims.UserID)
	assert.Equal(t, "127.0.0.1", claims.IP)
}

func TestGenerateRefreshToken(t *testing.T) {
	tokenService := NewTokenService("test_secret")
	refreshToken, hashedToken, err := tokenService.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)
	assert.NotEmpty(t, hashedToken)
}
