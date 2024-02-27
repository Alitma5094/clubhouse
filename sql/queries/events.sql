-- name: CreateEvent :one
INSERT INTO events (id,
                    created_at,
                    updated_at,
                    title,
                    start_at,
                    end_at,
                    user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
-- name: GetEvents :many
SELECT *
FROM events
WHERE user_id = $1;

-- name: DeleteEvent :exec
DELETE
FROM events
WHERE id = $1;