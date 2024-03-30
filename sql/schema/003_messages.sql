-- +goose Up
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    text TEXT NOT NULL
);
-- +goose Down
DROP TABLE messages;