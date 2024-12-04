CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id VARCHAR(36) PRIMARY KEY,
    hashed_token TEXT NOT NULL
);