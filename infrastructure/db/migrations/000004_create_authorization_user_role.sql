-- +goose Up
CREATE TABLE IF NOT EXISTS authz_user_role
(
    user_id   UUID NOT NULL REFERENCES authz_user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id   UUID NOT NULL REFERENCES authz_role (id) ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- +goose Down
DROP TABLE IF EXISTS authz_user_role;
