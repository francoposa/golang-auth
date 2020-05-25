-- +goose Up
CREATE TABLE IF NOT EXISTS resource
(
    id        TEXT NOT NULL PRIMARY KEY,
    resource  TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS resource;
