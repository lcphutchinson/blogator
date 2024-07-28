-- +goose Up
CREATE TABLE feed_follows(
	id uuid DEFAULT gen_random_uuid(),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	feed_id uuid NOT NULL REFERENCES feeds ON DELETE CASCADE,
	user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,

	PRIMARY KEY (id)
);

-- +goose StatementBegin
CREATE FUNCTION trg_follow() RETURNS trigger AS $$
	BEGIN
	INSERT INTO feed_follows(feed_id, user_id)
	VALUES (NEW.id, NEW.user_id);
	RETURN NULL;
	END;
	$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER follow_new_feed 
	AFTER INSERT ON feeds
	FOR EACH ROW
	EXECUTE FUNCTION trg_follow();


-- +goose Down
DROP FUNCTION trg_follow CASCADE;
DROP TABLE feed_follows;
