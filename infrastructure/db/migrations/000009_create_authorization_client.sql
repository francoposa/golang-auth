-- +goose Up
CREATE TABLE IF NOT EXISTS authz_client
(
    id            TEXT NOT NULL PRIMARY KEY,
    secret        TEXT,
    redirect_uri  TEXT NOT NULL,
    public        BOOL NOT NULL,
    first_party   BOOL NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS authz_client;