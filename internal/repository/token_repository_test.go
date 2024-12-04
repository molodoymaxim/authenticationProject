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
		hashed_token TEXT NOT NULL
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

	err := repo.SaveRefreshToken(userID, hashedToken)
	assert.NoError(t, err)

	storedHashedToken, err := repo.GetRefreshTokenHash(userID)
	assert.NoError(t, err)
	assert.Equal(t, hashedToken, storedHashedToken)

	newHashedToken := "new-hashed-refresh-token"
	err = repo.UpdateRefreshToken(userID, newHashedToken)
	assert.NoError(t, err)

	updatedHashedToken, err := repo.GetRefreshTokenHash(userID)
	assert.NoError(t, err)
	assert.Equal(t, newHashedToken, updatedHashedToken)
}
