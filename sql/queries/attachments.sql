-- name: CreateAttachment :one
INSERT INTO attachments (id, created_at, updated_at, media_type, url, message_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAttachments :many
SELECT *
FROM attachments
WHERE message_id = $1;

-- name: DeleteAttachment :exec
DELETE
FROM attachments
WHERE id = $1;
