-- +goose Up
CREATE TABLE users(
	id uuid DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name VARCHAR(32) NOT NULL,

	PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
