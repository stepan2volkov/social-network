-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id              BIGINT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    username        CHAR(150) NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    firstname       TEXT NOT NULL,
    lastname        TEXT NOT NULL,
    birthdate       DATE NOT NULL,
    gender          CHAR(1) NOT NULL,
    city_id         BIGINT NOT NULL,
    CONSTRAINT user_city_id_fk FOREIGN KEY (city_id) REFERENCES cities(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd