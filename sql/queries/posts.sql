-- name: CreatePost :one
INSERT INTO posts (title, url, description, feed_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetPostsByUser :many
SELECT *
FROM posts
WHERE feed_id IN (
	SELECT id
	FROM feeds
	WHERE user_id = $1
)
LIMIT $2;
