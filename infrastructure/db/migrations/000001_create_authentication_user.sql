-- +goose Up
CREATE TABLE IF NOT EXISTS authn_user
(
    id        UUID NOT NULL PRIMARY KEY,
    username  TEXT NOT NULL UNIQUE,
    email     TEXT NOT NULL,
    password  TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authn_user;
