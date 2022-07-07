-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers (
    follower_id         BIGINT NOT NULL,
    leader_id           BIGINT NOT NULL,
    CONSTRAINT followers_pk PRIMARY KEY (follower_id, leader_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE followers;
-- +goose StatementEnd
