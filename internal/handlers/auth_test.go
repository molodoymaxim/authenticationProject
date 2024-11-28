package handlers

import (
	"authenticationProject/internal/models"
	"authenticationProject/internal/repository"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockTokenRepository struct{}

func (m *MockTokenRepository) SaveRefreshToken(userID, hashedToken, accessToken string) error {
	return nil
}

func (m *MockTokenRepository) GetRefreshToken(userID string) (*models.RefreshTokenData, error) {
	return &models.RefreshTokenData{
		UserID:      userID,
		HashedToken: "$2a$10$examplehashedtoken",
		AccessToken: "exampleaccesstoken",
	}, nil
}

func (m *MockTokenRepository) UpdateRefreshToken(userID, hashedToken, accessToken string) error {
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
	// Инициализация зависимостей
	logger := logrus.New()
	tokenService := services.NewTokenService("test_secret")
	emailService := utils.NewEmailService("test_api_key")
	tokenRepository := repository.NewInMemoryTokenRepository() // Используйте in-memory репозиторий для тестов
	authHandler := NewAuthHandler(tokenService, tokenRepository, emailService, logger)

	// Генерация первоначальных токенов
	rr := httptest.NewRecorder()
	reqBody := `{"user_id": "test-user-id"}`
	req := httptest.NewRequest("POST", "/auth/token", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	authHandler.GenerateTokens(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", rr.Code)
	}

	var tokens models.TokenResponse
	err := json.Unmarshal(rr.Body.Bytes(), &tokens)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Подготовка запроса на обновление токенов
	rr = httptest.NewRecorder()
	reqBody = fmt.Sprintf(`{"refresh_token": "%s"}`, tokens.RefreshToken)
	req = httptest.NewRequest("POST", "/auth/refresh", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

	// Выполнение запроса на обновление токенов
	authHandler.RefreshTokens(rr, req)

	// Проверка результатов
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var newTokens models.TokenResponse
	err = json.Unmarshal(rr.Body.Bytes(), &newTokens)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("Access token should not be empty")
	}
	if newTokens.RefreshToken == "" {
		t.Error("Refresh token should not be empty")
	}
}
