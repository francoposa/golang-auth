-- +goose Up
CREATE TABLE IF NOT EXISTS client
(
    id        TEXT  NOT NULL PRIMARY KEY,
    secret    TEXT  NOT NULL,
    domain    TEXT  NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS client;