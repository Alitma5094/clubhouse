// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
        id,
        created_at,
        updated_at,
        email,
        name,
        hashed_password
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, email, name, hashed_password
`

type CreateUserParams struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	HashedPassword string    `json:"hashed_password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Email,
		arg.Name,
		arg.HashedPassword,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
		&i.HashedPassword,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at, email, name, hashed_password
FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
		&i.HashedPassword,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, name, hashed_password
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Name,
		&i.HashedPassword,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, created_at, updated_at, email, name, hashed_password
FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.Name,
			&i.HashedPassword,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET updated_at = $2,
    email = $3,
    name = $4,
    hashed_password = $5
WHERE id = $1
`

type UpdateUserParams struct {
	ID             uuid.UUID `json:"id"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	HashedPassword string    `json:"hashed_password"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.UpdatedAt,
		arg.Email,
		arg.Name,
		arg.HashedPassword,
	)
	return err
}
