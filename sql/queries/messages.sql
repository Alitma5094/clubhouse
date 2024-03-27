-- name: CreateMessage :one
INSERT INTO messages (
        id,
        created_at,
        updated_at,
        user_id,
        text,
        thread_id,
        attachments
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: GetMessages :many
SELECT *
FROM messages
WHERE thread_id = $1
ORDER BY created_at DESC;