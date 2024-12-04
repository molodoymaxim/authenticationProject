package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type TokenRepository interface {
	SaveRefreshToken(userID, hashedToken string) error
	GetRefreshTokenHash(userID string) (string, error)
	UpdateRefreshToken(userID, hashedToken string) error
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (r *tokenRepository) SaveRefreshToken(userID, hashedToken string) error {
	_, err := r.db.Exec(`
        INSERT INTO refresh_tokens (user_id, hashed_token)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO UPDATE
        SET hashed_token = EXCLUDED.hashed_token
    `, userID, hashedToken)
	return err
}

func (r *tokenRepository) GetRefreshTokenHash(userID string) (string, error) {
	var hashedToken string
	err := r.db.QueryRow(`
        SELECT hashed_token
        FROM refresh_tokens
        WHERE user_id = $1
    `, userID).Scan(&hashedToken)
	if err != nil {
		return "", err
	}
	return hashedToken, nil
}

func (r *tokenRepository) UpdateRefreshToken(userID, hashedToken string) error {
	_, err := r.db.Exec(`
        UPDATE refresh_tokens
        SET hashed_token = $1
        WHERE user_id = $2
    `, hashedToken, userID)
	return err
}
