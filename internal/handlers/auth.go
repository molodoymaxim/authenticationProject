package handlers

import (
	"authenticationProject/internal/repository"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	TokenService    *services.TokenService
	TokenRepository repository.TokenRepository
	EmailService    *utils.EmailService
	Logger          *logrus.Logger
}

func NewAuthHandler(
	tokenService *services.TokenService,
	tokenRepository repository.TokenRepository,
	emailService *utils.EmailService,
	logger *logrus.Logger,
) *AuthHandler {
	return &AuthHandler{
		TokenService:    tokenService,
		TokenRepository: tokenRepository,
		EmailService:    emailService,
		Logger:          logger,
	}
}

// GenerateTokens выдает пару Access и Refresh токенов
func (h *AuthHandler) GenerateTokens(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug("Received request to generate tokens")

	var req struct {
		UserID string `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.Warn("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	clientIP := r.RemoteAddr
	h.Logger.Infof("Generating tokens for user %s from IP %s", req.UserID, clientIP)

	accessToken, err := h.TokenService.GenerateAccessToken(req.UserID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate access token: ", err)
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, hashedToken, err := h.TokenService.GenerateRefreshToken()
	if err != nil {
		h.Logger.Error("Failed to generate refresh token: ", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = h.TokenRepository.SaveRefreshToken(req.UserID, hashedToken, accessToken)
	if err != nil {
		h.Logger.Error("Failed to save refresh token: ", err)
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Tokens generated successfully for user %s", req.UserID)
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	json.NewEncoder(w).Encode(response)
}

// RefreshTokens обновляет пару Access и Refresh токенов
func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	h.Logger.Debug("Received request to refresh tokens")

	authHeader := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.Logger.Warn("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := h.TokenService.ValidateAccessToken(accessToken)
	if err != nil {
		h.Logger.Warn("Invalid access token: ", err)
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	h.Logger.Infof("Refreshing tokens for user %s", userID)

	tokenData, err := h.TokenRepository.GetRefreshToken(userID)
	if err != nil {
		h.Logger.Warn("Refresh token not found for user: ", userID)
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	// Проверяем соответствие Refresh токена
	err = bcrypt.CompareHashAndPassword([]byte(tokenData.HashedToken), []byte(req.RefreshToken))
	if err != nil {
		h.Logger.Warn("Invalid refresh token for user: ", userID)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Проверяем изменение IP адреса
	clientIP := r.RemoteAddr
	if claims.IP != clientIP {
		h.Logger.Warnf("IP address changed for user %s: %s -> %s", userID, claims.IP, clientIP)
		h.EmailService.SendEmailWarning(userID, "IP address changed during token refresh.")
	}

	// Генерируем новые токены
	newAccessToken, err := h.TokenService.GenerateAccessToken(userID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate new access token: ", err)
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, newHashedToken, err := h.TokenService.GenerateRefreshToken()
	if err != nil {
		h.Logger.Error("Failed to generate new refresh token: ", err)
		http.Error(w, "Failed to generate new refresh token", http.StatusInternalServerError)
		return
	}

	// Обновляем токены в базе данных
	err = h.TokenRepository.UpdateRefreshToken(userID, newHashedToken, newAccessToken)
	if err != nil {
		h.Logger.Error("Failed to update refresh token: ", err)
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Tokens refreshed successfully for user %s", userID)
	response := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}
	json.NewEncoder(w).Encode(response)
}