// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feeds.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1, $2, $3)
RETURNING id, created_at, updated_at, name, url, user_id, last_fetched_at
`

type CreateFeedParams struct {
	Name   string    `json:"name"`
	Url    string    `json:"url"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed, arg.Name, arg.Url, arg.UserID)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const listFeeds = `-- name: ListFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
ORDER BY name
`

func (q *Queries) ListFeeds(ctx context.Context) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, listFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
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
