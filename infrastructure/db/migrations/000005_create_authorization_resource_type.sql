-- +goose Up
CREATE TABLE IF NOT EXISTS authz_resource_type
(
    id           UUID NOT NULL PRIMARY KEY,
    name         TEXT NOT NULL UNIQUE,
    description  TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authz_resource_type;
