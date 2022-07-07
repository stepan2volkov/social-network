-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cities (
    id      BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    name    TEXT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cities;
-- +goose StatementEnd
