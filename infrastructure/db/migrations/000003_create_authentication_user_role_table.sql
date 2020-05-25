-- +goose Up
CREATE TABLE IF NOT EXISTS authentication_user_role
(
    id       TEXT NOT NULL PRIMARY KEY,
    user_id  TEXT NOT NULL REFERENCES authentication_user (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role_id  TEXT NOT NULL REFERENCES authentication_role (id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS authentication_user_role;
