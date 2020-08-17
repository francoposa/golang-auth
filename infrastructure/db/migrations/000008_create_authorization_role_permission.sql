-- +goose Up
CREATE TABLE IF NOT EXISTS authz_role_permission
(
    role_id        UUID NOT NULL REFERENCES authz_role (id) ON UPDATE CASCADE ON DELETE CASCADE,
    permission_id  UUID NOT NULL REFERENCES authz_permission (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- +goose Down
DROP TABLE IF EXISTS authz_role_permission;