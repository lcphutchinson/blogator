-- +goose Up
CREATE TABLE posts(
	id uuid DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	title VARCHAR NOT NULL,
	url VARCHAR UNIQUE NOT NULL,
	description VARCHAR,
	feed_id uuid NOT NULL REFERENCES feeds ON DELETE CASCADE,

	PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE posts;
