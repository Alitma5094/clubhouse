// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: attachments.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createAttachment = `-- name: CreateAttachment :one
INSERT INTO attachments (id, created_at, updated_at, media_type, url, message_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, media_type, url, message_id
`

type CreateAttachmentParams struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	MediaType Media     `json:"media_type"`
	Url       string    `json:"url"`
	MessageID uuid.UUID `json:"message_id"`
}

func (q *Queries) CreateAttachment(ctx context.Context, arg CreateAttachmentParams) (Attachment, error) {
	row := q.db.QueryRowContext(ctx, createAttachment,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.MediaType,
		arg.Url,
		arg.MessageID,
	)
	var i Attachment
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.MediaType,
		&i.Url,
		&i.MessageID,
	)
	return i, err
}
