package repository

import (
	"authenticationProject/internal/models"
	"database/sql"

	_ "github.com/lib/pq"
)

type TokenRepository interface {
	SaveRefreshToken(userID, hashedToken, accessToken string) error
	GetRefreshToken(userID string) (*models.RefreshTokenData, error)
	UpdateRefreshToken(userID, hashedToken, accessToken string) error
}

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func NewDatabase(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// SaveRefreshToken сохраняет Refresh токен в базе данных
func (r *tokenRepository) SaveRefreshToken(userID, hashedToken, accessToken string) error {
	_, err := r.db.Exec(`
        INSERT INTO refresh_tokens (user_id, hashed_token, access_token)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE
        SET hashed_token = EXCLUDED.hashed_token,
            access_token = EXCLUDED.access_token
    `, userID, hashedToken, accessToken)
	return err
}

// GetRefreshToken получает Refresh токен из базы данных
func (r *tokenRepository) GetRefreshToken(userID string) (*models.RefreshTokenData, error) {
	var tokenData models.RefreshTokenData
	err := r.db.QueryRow(`
        SELECT user_id, hashed_token, access_token
        FROM refresh_tokens
        WHERE user_id = $1
    `, userID).Scan(&tokenData.UserID, &tokenData.HashedToken, &tokenData.AccessToken)
	if err != nil {
		return nil, err
	}
	return &tokenData, nil
}

// UpdateRefreshToken обновляет Refresh токен в базе данных
func (r *tokenRepository) UpdateRefreshToken(userID, hashedToken, accessToken string) error {
	_, err := r.db.Exec(`
        UPDATE refresh_tokens
        SET hashed_token = $1, access_token = $2
        WHERE user_id = $3
    `, hashedToken, accessToken, userID)
	return err
}
