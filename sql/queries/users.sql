-- name: CreateUser :one
INSERT INTO users (
        id,
        created_at,
        updated_at,
        email,
        name,
        hashed_password
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: UpdateUser :exec
UPDATE users
SET updated_at = $2,
    email = $3,
    name = $4,
    hashed_password = $5
WHERE id = $1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;
-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;
-- name: GetUsers :many
SELECT *
FROM users;