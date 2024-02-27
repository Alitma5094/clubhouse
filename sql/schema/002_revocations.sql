-- +goose Up
CREATE TABLE revocations (
    id UUID PRIMARY KEY,
    revoked_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    token TEXT UNIQUE NOT NULL,
    user_id UUID REFERENCES users (id) NOT NULL
);
-- +goose Down
DROP TABLE revocations;