-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profiles (
	id				BINARY(16) NOT NULL PRIMARY KEY,
	username		VARCHAR(200) NOT NULL UNIQUE,
	password_hash	VARCHAR(200) NOT NULL,
	firstname		VARCHAR(100) NOT NULL,
	lastname		VARCHAR(100) NOT NULL,
	birthdate		DATE NOT NULL,
	biography		TEXT NOT NULL,
	city			VARCHAR(100) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS profiles;
-- +goose StatementEnd
