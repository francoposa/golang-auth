-- +goose Up
CREATE TABLE IF NOT EXISTS authz_user
(
    id        UUID NOT NULL PRIMARY KEY,
    username  TEXT NOT NULL UNIQUE,
    email     TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS authz_user;
