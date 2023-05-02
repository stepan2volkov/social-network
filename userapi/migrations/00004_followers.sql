-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers (
	follower	VARCHAR(200) NOT NULL REFERENCES profiles(username),
	followee	VARCHAR(200) NOT NULL REFERENCES profiles(username),
	CONSTRAINT followers_check_names_are_not_equal CHECK (followee != follower),
	CONSTRAINT followers_unique_row UNIQUE (follower, followee)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers;
-- +goose StatementEnd
