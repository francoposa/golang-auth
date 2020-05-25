-- +goose Up
CREATE TABLE IF NOT EXISTS resource
(
    id   TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS resource;
