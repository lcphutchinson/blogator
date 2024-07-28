-- name: CreateFollow :one
INSERT INTO feed_follows (feed_id, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveFollow :execresult 
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;

-- name: ListUserFollows :many
SELECT *
FROM feed_follows
WHERE user_id = $1
ORDER BY updated_at DESC;
