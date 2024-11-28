package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open("postgres", "postgres://postgres:148822@localhost:5432/authdb_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	_, err = testDB.Exec(`CREATE TABLE IF NOT EXISTS refresh_tokens (
        user_id VARCHAR(36) PRIMARY KEY,
        hashed_token TEXT NOT NULL,
        access_token TEXT NOT NULL
    );`)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	testDB.Exec(`DROP TABLE refresh_tokens;`)

	os.Exit(code)
}

func TestTokenRepository(t *testing.T) {
	repo := NewTokenRepository(testDB)

	userID := "test-user-id"
	hashedToken := "hashed-refresh-token"
	accessToken := "access-token"

	err := repo.SaveRefreshToken(userID, hashedToken, accessToken)
	assert.NoError(t, err)

	tokenData, err := repo.GetRefreshToken(userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, tokenData.UserID)
	assert.Equal(t, hashedToken, tokenData.HashedToken)
	assert.Equal(t, accessToken, tokenData.AccessToken)

	newHashedToken := "new-hashed-refresh-token"
	newAccessToken := "new-access-token"
	err = repo.UpdateRefreshToken(userID, newHashedToken, newAccessToken)
	assert.NoError(t, err)

	updatedTokenData, err := repo.GetRefreshToken(userID)
	assert.NoError(t, err)
	assert.Equal(t, newHashedToken, updatedTokenData.HashedToken)
	assert.Equal(t, newAccessToken, updatedTokenData.AccessToken)
}
