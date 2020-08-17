-- +goose Up
CREATE TABLE IF NOT EXISTS authz_role
(
    id    UUID NOT NULL PRIMARY KEY,
    name  TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS authz_role;
