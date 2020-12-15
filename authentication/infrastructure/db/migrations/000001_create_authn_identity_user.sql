-- +goose Up
CREATE TABLE IF NOT EXISTS authn_identity_user
(
    id UUID NOT NULL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authn_identity_user;
