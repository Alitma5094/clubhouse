-- name: CreateThread :one
INSERT INTO threads (
        id,
        created_at,
        updated_at,
        user_id,
        title
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: GetThreads :many
SELECT threads.*
FROM threads
    INNER JOIN users_threads ON threads.id = users_threads.thread_id
WHERE users_threads.user_id = $1;
-- name: SubscribeToThread :one
INSERT INTO users_threads (id, created_at, updated_at, user_id, thread_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: DeleteThread :exec
DELETE FROM users_threads
WHERE thread_id = $1;
DELETE FROM threads
WHERE id = $1;
-- name: GetUserSubscribedThreads :many
SELECT thread_id
FROM users_threads
WHERE user_id = $1;
-- name: GetSubscribedUsers :many
SELECT users.*
FROM users
    INNER JOIN users_threads ON users.id = users_threads.user_id
WHERE users_threads.thread_id = $1;
-- name: GetUnsubscribedUsers :many
SELECT users.*
FROM users
WHERE users.id NOT IN (
        SELECT user_id
        FROM users_threads
        WHERE thread_id = $1
    );
-- name: UnsubscribeFromThread :exec
DELETE FROM users_threads
WHERE user_id = $1
    AND thread_id = $2;