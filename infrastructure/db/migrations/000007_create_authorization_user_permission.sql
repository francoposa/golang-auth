-- +goose Up
CREATE TABLE IF NOT EXISTS authz_user_permission
(
    user_id        UUID NOT NULL REFERENCES authz_user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    permission_id  UUID NOT NULL REFERENCES authz_permission (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

-- +goose Down
DROP TABLE IF EXISTS authz_user_permission;