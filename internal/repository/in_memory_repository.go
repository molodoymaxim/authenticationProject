package repository

import (
	"errors"
	"sync"
)

type InMemoryTokenRepository struct {
	tokens map[string]string
	mu     sync.Mutex
}

var (
	ErrTokenNotFound = errors.New("Token not found")
)

func NewInMemoryTokenRepository() TokenRepository {
	return &InMemoryTokenRepository{
		tokens: make(map[string]string),
	}
}

func (r *InMemoryTokenRepository) SaveRefreshToken(userID, hashedToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[userID] = hashedToken
	return nil
}

func (r *InMemoryTokenRepository) GetRefreshTokenHash(userID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	hashedToken, exists := r.tokens[userID]
	if !exists {
		return "", ErrTokenNotFound
	}
	return hashedToken, nil
}

func (r *InMemoryTokenRepository) UpdateRefreshToken(userID, hashedToken string) error {
	return r.SaveRefreshToken(userID, hashedToken)
}
