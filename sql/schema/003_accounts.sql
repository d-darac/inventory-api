-- +goose Up
CREATE TABLE accounts 
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    country    countries NOT NULL,
    deleted    BOOLEAN NOT NULL DEFAULT FALSE,
    nickname   TEXT,
    owner_id   UUID DEFAULT NULL REFERENCES users ON DELETE SET NULL
);

-- +goose Down
DROP TABLE accounts;