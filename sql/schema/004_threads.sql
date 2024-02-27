-- +goose Up
CREATE TABLE threads (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users (id) NOT NULL,
    title TEXT NOT NULL
);
ALTER TABLE messages
ADD COLUMN thread_id UUID REFERENCES threads (id) NOT NULL;
CREATE TABLE users_threads (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users (id) NOT NULL,
    thread_id UUID REFERENCES threads (id) NOT NULL,
    UNIQUE (user_id, thread_id)
);
-- +goose Down
ALTER TABLE messages DROP COLUMN thread_id;
DROP TABLE threads;
DROP TABLE users_threads;