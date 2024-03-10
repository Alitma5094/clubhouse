// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: threads.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createThread = `-- name: CreateThread :one
INSERT INTO threads (
        id,
        created_at,
        updated_at,
        user_id,
        title
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, title
`

type CreateThreadParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
}

func (q *Queries) CreateThread(ctx context.Context, arg CreateThreadParams) (Thread, error) {
	row := q.db.QueryRowContext(ctx, createThread,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.Title,
	)
	var i Thread
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Title,
	)
	return i, err
}

const deleteThread = `-- name: DeleteThread :exec
DELETE FROM users_threads
WHERE thread_id = $1
`

func (q *Queries) DeleteThread(ctx context.Context, threadID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteThread, threadID)
	return err
}

const getThreads = `-- name: GetThreads :many
SELECT threads.id, threads.created_at, threads.updated_at, threads.user_id, threads.title
FROM threads
    INNER JOIN users_threads ON threads.id = users_threads.thread_id
WHERE users_threads.user_id = $1
`

func (q *Queries) GetThreads(ctx context.Context, userID uuid.UUID) ([]Thread, error) {
	rows, err := q.db.QueryContext(ctx, getThreads, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Thread
	for rows.Next() {
		var i Thread
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Title,
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

const getUserSubscribedThreads = `-- name: GetUserSubscribedThreads :many
SELECT thread_id
FROM users_threads
WHERE user_id = $1
`

func (q *Queries) GetUserSubscribedThreads(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getUserSubscribedThreads, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var thread_id uuid.UUID
		if err := rows.Scan(&thread_id); err != nil {
			return nil, err
		}
		items = append(items, thread_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const subscribeToThread = `-- name: SubscribeToThread :one
INSERT INTO users_threads (id, created_at, updated_at, user_id, thread_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, thread_id
`

type SubscribeToThreadParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	ThreadID  uuid.UUID `json:"thread_id"`
}

func (q *Queries) SubscribeToThread(ctx context.Context, arg SubscribeToThreadParams) (UsersThread, error) {
	row := q.db.QueryRowContext(ctx, subscribeToThread,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.ThreadID,
	)
	var i UsersThread
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ThreadID,
	)
	return i, err
}
