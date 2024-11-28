package repository

import (
	"authenticationProject/internal/models"
	"errors"
	"sync"
)

type InMemoryTokenRepository struct {
	tokens map[string]*models.RefreshTokenData
	mu     sync.Mutex
}

var (
	ErrTokenNotFound = errors.New("Token not found")
)

func NewInMemoryTokenRepository() TokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]*models.RefreshTokenData),
	}
}

func (r *InMemoryTokenRepository) SaveRefreshToken(userID, hashedToken, accessToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[userID] = &models.RefreshTokenData{
		HashedToken: hashedToken,
		AccessToken: accessToken,
	}
	return nil
}

func (r *InMemoryTokenRepository) GetRefreshToken(userID string) (*models.RefreshTokenData, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tokenData, exists := r.tokens[userID]
	if !exists {
		return nil, ErrTokenNotFound
	}
	return tokenData, nil
}

func (r *InMemoryTokenRepository) UpdateRefreshToken(userID, hashedToken, accessToken string) error {
	return r.SaveRefreshToken(userID, hashedToken, accessToken)
}
