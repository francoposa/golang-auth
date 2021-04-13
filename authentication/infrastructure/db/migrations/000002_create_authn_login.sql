-- +goose Up
CREATE TABLE IF NOT EXISTS authn_login
(
    id UUID NOT NULL,
    redirect_url TEXT NOT NULL,
    status TEXT NOT NULL,
    attempts INTEGER NOT NULL,
    csrf_token TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS authn_login;
