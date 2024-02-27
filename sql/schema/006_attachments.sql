-- +goose Up
CREATE TYPE media AS ENUM ('image', 'video', 'document');
CREATE TABLE attachments (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    media_type media NOT NULL,
    url TEXT NOT NULL,
    message_id UUID REFERENCES messages (id) NOT NULL
);
-- +goose Down
DROP TABLE attachments;
DROP TYPE media;