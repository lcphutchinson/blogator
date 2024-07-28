-- +goose Up
CREATE TABLE feeds(
	id uuid DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name VARCHAR(32) NOT NULL,
	url VARCHAR UNIQUE NOT NULL,
	user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,

	PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE feeds;
