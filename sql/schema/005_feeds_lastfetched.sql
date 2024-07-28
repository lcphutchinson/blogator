-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;

-- +goose StatementBegin
CREATE FUNCTION fetch_n_feeds(INTEGER) RETURNS SETOF feeds AS $$
	BEGIN
	RETURN QUERY
	UPDATE feeds
	SET last_fetched_at = NOW(),
	updated_at = NOW()
	WHERE id IN (
		SELECT id
		FROM feeds
		ORDER BY last_fetched_at ASC
		NULLS FIRST
		LIMIT $1
	)
	RETURNING *;
	END;
	$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION fetch_n_feeds;
ALTER TABLE feeds
DROP COLUMN last_fetched_at;
