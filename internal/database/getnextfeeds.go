package database

import "context"


const getNextNFeeds = `-- name: GetNextNFeeds :many
SELECT * FROM fetch_n_feeds($1)
`

func (q *Queries) GetNextNFeeds(ctx context.Context, fetchNFeeds int32) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, getNextNFeeds, fetchNFeeds)
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
