package handlers

import (
	"authenticationProject/internal/repository"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

// GenerateTokens генерирует пару Access и Refresh токенов
// @Summary Генерация токенов
// @Description Генерирует Access и Refresh токены для указанного пользователя
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body models.TokenRequest true "User ID"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/token [post]
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

	userID := req.UserID
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	h.Logger.Infof("Generating tokens for user %s from IP %s", userID, clientIP)

	pairID := uuid.NewString()

	accessToken, err := h.TokenService.GenerateAccessToken(userID, pairID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate access token: ", err)
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshTokenEncoded, hashedToken, err := h.TokenService.GenerateRefreshToken(userID, pairID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate refresh token: ", err)
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	err = h.TokenRepository.SaveRefreshToken(userID, hashedToken)
	if err != nil {
		h.Logger.Error("Failed to save refresh token: ", err)
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Tokens generated successfully for user %s", userID)
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshTokenEncoded,
	}
	json.NewEncoder(w).Encode(response)
}

// RefreshTokens обновляет пару токенов
// @Summary Обновление токенов
// @Description Обновляет Access и Refresh токены
// @Tags Authentication
// @Accept json
// @Produce json
// @Param data body models.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Security BearerAuth
// @Router /auth/refresh [post]
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
	accessPairID := claims.PairID
	h.Logger.Infof("Refreshing tokens for user %s", userID)

	storedHash, err := h.TokenRepository.GetRefreshTokenHash(userID)
	if err != nil {
		h.Logger.Warn("Refresh token not found for user: ", userID)
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	refreshTokenBytes, err := base64.StdEncoding.DecodeString(req.RefreshToken)
	if err != nil {
		h.Logger.Warn("Invalid refresh token format")
		http.Error(w, "Invalid refresh token format", http.StatusBadRequest)
		return
	}
	refreshTokenPlain := string(refreshTokenBytes)

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refreshTokenPlain))
	if err != nil {
		h.Logger.Warn("Invalid refresh token for user: ", userID)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(refreshTokenPlain, ":")
	if len(parts) != 3 {
		h.Logger.Warn("Invalid refresh token format")
		http.Error(w, "Invalid refresh token format", http.StatusBadRequest)
		return
	}
	refreshUserID := parts[0]
	refreshPairID := parts[1]
	refreshIP := parts[2]

	if userID != refreshUserID || accessPairID != refreshPairID {
		h.Logger.Warn("Access and Refresh tokens do not match")
		http.Error(w, "Invalid token pair", http.StatusUnauthorized)
		return
	}

	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if clientIP != refreshIP {
		h.Logger.Warnf("IP address changed for user %s: %s -> %s", userID, refreshIP, clientIP)
		h.EmailService.SendEmailWarning(userID, "IP address changed during token refresh.")
	}

	newPairID := uuid.NewString()
	newAccessToken, err := h.TokenService.GenerateAccessToken(userID, newPairID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate new access token: ", err)
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	newRefreshTokenEncoded, newHashedToken, err := h.TokenService.GenerateRefreshToken(userID, newPairID, clientIP)
	if err != nil {
		h.Logger.Error("Failed to generate new refresh token: ", err)
		http.Error(w, "Failed to generate new refresh token", http.StatusInternalServerError)
		return
	}

	err = h.TokenRepository.UpdateRefreshToken(userID, newHashedToken)
	if err != nil {
		h.Logger.Error("Failed to update refresh token: ", err)
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	h.Logger.Infof("Tokens refreshed successfully for user %s", userID)
	response := map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshTokenEncoded,
	}
	json.NewEncoder(w).Encode(response)
}

// ProtectedHandler пример защищенного маршрута
// @Summary Пример защищенного маршрута
// @Description Доступен только для авторизованных пользователей
// @Tags Protected
// @Produce plain
// @Success 200 {string} string "Hello, user"
// @Failure 401 {string} string "Unauthorized"
// @Security BearerAuth
// @Router /protected [get]
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Hello, user " + userID))
}
