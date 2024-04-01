-- name: CreateFcmToken :one
INSERT INTO fcm_tokens (
        id,
        user_id,
        token,
        created_at
    )
VALUES ($1, $2, $3, $4)
RETURNING *;
-- name: GetFcmTokens :many
SELECT fcm_tokens.*
FROM fcm_tokens
    JOIN users_threads ON fcm_tokens.user_id = users_threads.user_id
WHERE users_threads.thread_id = $1;