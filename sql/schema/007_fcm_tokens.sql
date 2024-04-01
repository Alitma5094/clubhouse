-- +goose Up
CREATE TABLE fcm_tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose Down
DROP TABLE fcm_tokens;