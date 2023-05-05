-- +goose Up
-- +goose StatementBegin
CREATE INDEX profiles_name_search_idx ON profiles (firstname, lastname);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX profiles_name_search_idx ON profiles;
-- +goose StatementEnd
