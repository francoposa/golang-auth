-- +goose Up
CREATE TABLE IF NOT EXISTS authz_permission
(
    id                UUID NOT NULL PRIMARY KEY,
    resource_type_id  UUID NOT NULL REFERENCES authz_resource_type (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name              TEXT NOT NULL UNIQUE,
    description       TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authz_permission;