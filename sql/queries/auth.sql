-- name: GetRevokedToken :one
SELECT *
FROM revocations
WHERE token = $1;
-- name: CreateRevocation :one
INSERT INTO revocations (id, revoked_at, updated_at, token, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;