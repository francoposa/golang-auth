-- +goose Up
CREATE TABLE IF NOT EXISTS auth_user
(
    id       TEXT NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email    TEXT NOT NULL,
    password TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS auth_user;
