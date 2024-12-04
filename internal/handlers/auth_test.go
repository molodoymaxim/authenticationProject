package handlers

import (
	"authenticationProject/internal/repository"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockTokenRepository struct{}

func (m *MockTokenRepository) SaveRefreshToken(userID, hashedToken string) error {
	return nil
}

func (m *MockTokenRepository) GetRefreshTokenHash(userID string) (string, error) {
	return "$2a$10$examplehashedtoken", nil
}

func (m *MockTokenRepository) UpdateRefreshToken(userID, hashedToken string) error {
	return nil
}

func TestGenerateTokens(t *testing.T) {
	tokenService := services.NewTokenService("test_secret")
	tokenRepository := &MockTokenRepository{}
	emailService := utils.NewEmailService("test_api_key")
	logger := logrus.New()

	authHandler := NewAuthHandler(tokenService, tokenRepository, emailService, logger)

	r := chi.NewRouter()
	r.Post("/auth/token", authHandler.GenerateTokens)

	reqBody := map[string]string{
		"user_id": "test-user-id",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "/auth/token", bytes.NewBuffer(bodyBytes))
	assert.NoError(t, err)
	req.RemoteAddr = "127.0.0.1:12345"

	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	var respBody map[string]string
	err = json.Unmarshal(recorder.Body.Bytes(), &respBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, respBody["access_token"])
	assert.NotEmpty(t, respBody["refresh_token"])
}

func TestRefreshTokens(t *testing.T) {
	logger := logrus.New()
	tokenService := services.NewTokenService("test_secret")
	emailService := utils.NewEmailService("test_api_key")
	tokenRepository := repository.NewInMemoryTokenRepository()
	authHandler := NewAuthHandler(tokenService, tokenRepository, emailService, logger)

	rr := httptest.NewRecorder()
	reqBody := `{"user_id": "test-user-id"}`
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345"

	authHandler.GenerateTokens(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tokens map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &tokens)
	assert.NoError(t, err)
	accessToken := tokens["access_token"]
	refreshToken := tokens["refresh_token"]

	rr = httptest.NewRecorder()
	reqBody = `{"refresh_token": "` + refreshToken + `"}`
	req = httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.RemoteAddr = "127.0.0.1:12345"

	authHandler.RefreshTokens(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var newTokens map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &newTokens)
	assert.NoError(t, err)

	assert.NotEmpty(t, newTokens["access_token"])
	assert.NotEmpty(t, newTokens["refresh_token"])
}
