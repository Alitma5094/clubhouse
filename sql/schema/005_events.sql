-- +goose Up
CREATE TABLE events (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    user_id UUID REFERENCES users (id) NOT NULL
);
-- +goose Down
DROP TABLE events;