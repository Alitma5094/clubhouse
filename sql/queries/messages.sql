-- name: CreateMessage :one
INSERT INTO messages (
        id,
        created_at,
        updated_at,
        user_id,
        text,
        thread_id
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetMessages :many
SELECT *
FROM messages
WHERE thread_id = $1
ORDER BY created_at DESC;
-- name: GetMessagesWithAttachment :many
SELECT m.id,
    m.created_at,
    m.updated_at,
    m.user_id,
    m.text,
    m.thread_id,
    a.media_type as attachment_media_type,
    a.url as attachment_url
FROM messages m
    LEFT JOIN attachments a ON m.id = a.message_id
WHERE m.thread_id = $1
ORDER BY m.created_at DESC;