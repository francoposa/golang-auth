-- +goose Up
CREATE TABLE IF NOT EXISTS authn_login
(
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES authn_user(id)
);

-- +goose Down
DROP TABLE IF EXISTS authn_login;