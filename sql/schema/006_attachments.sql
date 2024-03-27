-- +goose Up
ALTER TABLE messages
ADD COLUMN attachments TEXT [];
-- +goose Down
ALTER TABLE messages DROP COLUMN attachments;