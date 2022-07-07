
-- +goose Up
-- +goose StatementBegin
INSERT INTO 
    cities(name) 
VALUES 
    ('Москва'),
    ('Санкт-Питербург'),
    ('Челябинск');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE cities;
-- +goose StatementEnd